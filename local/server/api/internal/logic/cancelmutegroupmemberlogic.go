package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type CancelMuteGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMuteGroupMemberLogic {
	return &CancelMuteGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMuteGroupMemberLogic) CancelMuteGroupMember(req *types.CancelMuteGroupMemberReq) (resp *types.CancelMuteGroupMemberResp, err error) {
	rpcReq := &user.CancelMuteGroupMemberReq{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}

	_, err = l.svcCtx.GroupClient.CancelMuteGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.CancelMuteGroupMemberResp{
	}

	return resp, nil
}
