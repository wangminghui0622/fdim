package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetFullJoinGroupIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFullJoinGroupIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullJoinGroupIDsLogic {
	return &GetFullJoinGroupIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFullJoinGroupIDsLogic) GetFullJoinGroupIDs(req *types.GetFullJoinGroupIDsReq) (resp *types.GetFullJoinGroupIDsResp, err error) {
	rpcReq := &user.GetFullJoinGroupIDsReq{
		UserID: req.UserID,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetFullJoinGroupIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetFullJoinGroupIDsResp{
		GroupIDs: rpcResp.GroupIDs,
	}

	return resp, nil
}
