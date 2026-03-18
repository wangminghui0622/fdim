package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type InviteUserToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteUserToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteUserToGroupLogic {
	return &InviteUserToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteUserToGroupLogic) InviteUserToGroup(req *types.InviteUserToGroupReq) (resp *types.InviteUserToGroupResp, err error) {
	rpcReq := &user.InviteUserToGroupReq{
		GroupID:        req.GroupID,
		Reason:         req.Reason,
		InvitedUserIDs: req.InvitedUserIDs,
	}

	if req.SendMessage != nil {
		rpcReq.SendMessage = req.SendMessage
	}

	_, err = l.svcCtx.GroupClient.InviteUserToGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.InviteUserToGroupResp{
	}

	return resp, nil
}
