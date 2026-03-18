package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetSelfUnhandledApplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelfUnhandledApplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfUnhandledApplyCountLogic {
	return &GetSelfUnhandledApplyCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelfUnhandledApplyCountLogic) GetSelfUnhandledApplyCount(req *types.GetSelfUnhandledApplyCountReq) (resp *types.GetSelfUnhandledApplyCountResp, err error) {
	rpcReq := &user.GetSelfUnhandledApplyCountReq{
		UserID: req.UserID,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetSelfUnhandledApplyCount(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetSelfUnhandledApplyCountResp{
		Count: rpcResp.Count,
	}

	return resp, nil
}
