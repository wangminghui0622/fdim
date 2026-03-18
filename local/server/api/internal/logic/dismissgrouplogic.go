package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type DismissGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDismissGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupLogic {
	return &DismissGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DismissGroupLogic) DismissGroup(req *types.DismissGroupReq) (resp *types.DismissGroupResp, err error) {
	rpcReq := &user.DismissGroupReq{
		GroupID: req.GroupID,
	}

	_, err = l.svcCtx.GroupClient.DismissGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.DismissGroupResp{
	}

	return resp, nil
}
