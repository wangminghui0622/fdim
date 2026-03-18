package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupMemberUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberUserIDsLogic {
	return &GetGroupMemberUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMemberUserIDsLogic) GetGroupMemberUserIDs(req *types.GetGroupMemberUserIDsReq) (resp *types.GetGroupMemberUserIDsResp, err error) {
	rpcReq := &user.GetGroupMemberUserIDsReq{
		GroupID: req.GroupID,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupMemberUserIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetGroupMemberUserIDsResp{
		UserIDs: rpcResp.UserIDs,
	}

	return resp, nil
}
