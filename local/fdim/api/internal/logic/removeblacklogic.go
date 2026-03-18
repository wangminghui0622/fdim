package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type RemoveBlackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveBlackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBlackLogic {
	return &RemoveBlackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveBlackLogic) RemoveBlack(req *types.RemoveBlackReq) (resp *types.RemoveBlackResp, err error) {
	rpcReq := &user.RemoveBlackReq{
		OwnerUserID: req.OwnerUserID,
		BlackUserID: req.BlackUserID,
	}

	_, err = l.svcCtx.FriendClient.RemoveBlack(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.RemoveBlackResp{
	}

	return resp, nil
}
