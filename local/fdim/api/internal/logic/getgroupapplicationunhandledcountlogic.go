package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupApplicationUnhandledCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupApplicationUnhandledCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupApplicationUnhandledCountLogic {
	return &GetGroupApplicationUnhandledCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupApplicationUnhandledCountLogic) GetGroupApplicationUnhandledCount(req *types.GetGroupApplicationUnhandledCountReq) (resp *types.GetGroupApplicationUnhandledCountResp, err error) {
	rpcReq := &user.GetGroupApplicationUnhandledCountReq{
		UserID: req.UserID,
		Time:   req.Time,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupApplicationUnhandledCount(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetGroupApplicationUnhandledCountResp{
		Count: rpcResp.Count,
	}

	return resp, nil
}
