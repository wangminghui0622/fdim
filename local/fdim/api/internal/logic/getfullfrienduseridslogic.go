package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetFullFriendUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFullFriendUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullFriendUserIDsLogic {
	return &GetFullFriendUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFullFriendUserIDsLogic) GetFullFriendUserIDs(req *types.GetFullFriendUserIDsReq) (resp *types.GetFullFriendUserIDsResp, err error) {
	rpcReq := &user.GetFullFriendUserIDsReq{
		UserID: req.UserID,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetFullFriendUserIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetFullFriendUserIDsResp{
		UserIDs: rpcResp.UserIDs,
	}

	return resp, nil
}
