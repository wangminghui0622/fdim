package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MuteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMuteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupLogic {
	return &MuteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MuteGroupLogic) MuteGroup(req *user.MuteGroupReq) (*user.MuteGroupResp, error) {
	resp := &user.MuteGroupResp{}

	// ������֤
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// Ȩ�޼�飺��Ҫ��Ⱥ�������Ա
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
		if member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin {
			return nil, fmt.Errorf("no permission: only group owner or admin can mute group")
		}
	}

	// ����Ⱥ��״̬Ϊ����
	data := map[string]interface{}{
		"status": constant.GroupStatusMuted,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMap(l.ctx, req.GroupID, data); err != nil {
		return nil, err
	}

	// ����Ⱥ�����֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupMutedNotification(l.ctx, req.GroupID, opUserID)
	}

	return resp, nil
}
