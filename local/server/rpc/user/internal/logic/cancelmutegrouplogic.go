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

type CancelMuteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelMuteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMuteGroupLogic {
	return &CancelMuteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelMuteGroupLogic) CancelMuteGroup(req *user.CancelMuteGroupReq) (*user.CancelMuteGroupResp, error) {
	resp := &user.CancelMuteGroupResp{}

	// ???????
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// ??????????????????
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
		if member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin {
			return nil, fmt.Errorf("no permission: only group owner or admin can cancel mute group")
		}
	}

	// ??????????????
	data := map[string]interface{}{
		"status": constant.GroupStatusNormal,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMap(l.ctx, req.GroupID, data); err != nil {
		return nil, err
	}

	// ????????????????
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupCancelMutedNotification(l.ctx, req.GroupID, opUserID)
	}

	return resp, nil
}
