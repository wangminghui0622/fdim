package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetFullGroupMemberUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFullGroupMemberUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullGroupMemberUserIDsLogic {
	return &GetFullGroupMemberUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFullGroupMemberUserIDsLogic) GetFullGroupMemberUserIDs(req *types.GetFullGroupMemberUserIDsReq) (resp *types.GetFullGroupMemberUserIDsResp, err error) {
	rpcReq := &user.GetFullGroupMemberUserIDsReq{
		GroupID: req.GroupID,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetFullGroupMemberUserIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetFullGroupMemberUserIDsResp{
		UserIDs: rpcResp.UserIDs,
	}

	return resp, nil
}
