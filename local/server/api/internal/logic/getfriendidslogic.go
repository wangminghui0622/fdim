package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetFriendIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendIDsLogic {
	return &GetFriendIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendIDsLogic) GetFriendIDs(req *types.GetFriendIDsReq) (resp *types.GetFriendIDsResp, err error) {
	rpcReq := &user.GetFriendIDsReq{
		UserID: req.UserID,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetFriendIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.GetFriendIDsResp{
		FriendIDs: rpcResp.FriendIDs,
	}

	return resp, nil
}
