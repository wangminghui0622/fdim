package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type AddBlackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddBlackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBlackLogic {
	return &AddBlackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBlackLogic) AddBlack(req *types.AddBlackReq) (resp *types.AddBlackResp, err error) {
	rpcReq := &user.AddBlackReq{
		OwnerUserID: req.OwnerUserID,
		BlackUserID: req.BlackUserID,
	}

	_, err = l.svcCtx.FriendClient.AddBlack(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.AddBlackResp{
	}

	return resp, nil
}
