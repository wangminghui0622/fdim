package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type CancelMuteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMuteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMuteGroupLogic {
	return &CancelMuteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMuteGroupLogic) CancelMuteGroup(req *types.CancelMuteGroupReq) (resp *types.CancelMuteGroupResp, err error) {
	rpcReq := &user.CancelMuteGroupReq{
		GroupID: req.GroupID,
	}

	_, err = l.svcCtx.GroupClient.CancelMuteGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.CancelMuteGroupResp{
	}

	return resp, nil
}
