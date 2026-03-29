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

type CancelMuteGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMuteGroupMemberLogic {
	return &CancelMuteGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelMuteGroupMemberLogic) CancelMuteGroupMember(req *user.CancelMuteGroupMemberReq) (*user.CancelMuteGroupMemberResp, error) {
	resp := &user.CancelMuteGroupMemberResp{}

	// ???????
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// ????????????????????????????????????
	opUserID := mcontext.GetOpUserID(l.ctx)
	member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}

	if !authverify.IsAdmin(l.ctx) {
		// ???????????????
		if member.RoleLevel == constant.GroupOwner {
			return nil, fmt.Errorf("cannot cancel mute group owner")
		}

		opMember, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}

		// ??????
		if member.RoleLevel == constant.GroupAdmin {
			// ???????????????????????
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, fmt.Errorf("no permission: only group owner can cancel mute group admin")
			}
		} else if member.RoleLevel == constant.GroupOrdinaryUsers {
			// ???????????????????????????
			if opMember.RoleLevel != constant.GroupOwner && opMember.RoleLevel != constant.GroupAdmin {
				return nil, fmt.Errorf("no permission: only group owner or admin can cancel mute member")
			}
		}
	}

	// ????????????????0??????????
	data := map[string]interface{}{
		"mute_end_time": time.UnixMilli(0),
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMemberMap(l.ctx, req.GroupID, req.UserID, data); err != nil {
		return nil, err
	}

	// ?????????????????
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupMemberCancelMutedNotification(l.ctx, req.GroupID, opUserID, []string{req.UserID})
	}

	return resp, nil
}
