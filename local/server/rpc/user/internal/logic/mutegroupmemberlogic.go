package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MuteGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupMemberLogic {
	return &MuteGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MuteGroupMemberLogic) MuteGroupMember(req *user.MuteGroupMemberReq) (*user.MuteGroupMemberResp, error) {
	resp := &user.MuteGroupMemberResp{}

	// ������֤
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// Ȩ�޼�飺��Ҫ��Ⱥ�������Ա���Ҳ��ܽ���Ⱥ��
	opUserID := mcontext.GetOpUserID(l.ctx)
	member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}

	if !authverify.IsAdmin(l.ctx) {
		// ���ܽ���Ⱥ��
		if member.RoleLevel == constant.GroupOwner {
			return nil, fmt.Errorf("cannot mute group owner")
		}

		opMember, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}

		// ���Ȩ��
		if member.RoleLevel == constant.GroupAdmin {
			// ֻ��Ⱥ�����Խ��Թ���Ա
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, fmt.Errorf("no permission: only group owner can mute group admin")
			}
		} else if member.RoleLevel == constant.GroupOrdinaryUsers {
			// Ⱥ�������Ա���Խ�����ͨ��Ա
			if opMember.RoleLevel != constant.GroupOwner && opMember.RoleLevel != constant.GroupAdmin {
				return nil, fmt.Errorf("no permission: only group owner or admin can mute member")
			}
		}
	}

	// ������Խ���ʱ��
	muteEndTime := time.Now().Add(time.Duration(req.MutedSeconds) * time.Second)

	// ����Ⱥ��Ա����ʱ��
	data := map[string]interface{}{
		"mute_end_time": muteEndTime,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMemberMap(l.ctx, req.GroupID, req.UserID, data); err != nil {
		return nil, err
	}

	// ����Ⱥ��Ա����֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupMemberMutedNotification(l.ctx, req.GroupID, opUserID, []string{req.UserID})
	}

	return resp, nil
}
