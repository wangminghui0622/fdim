package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFriendLogic) DeleteFriend(req *types.DeleteFriendReq) (resp *types.DeleteFriendResp, err error) {
	rpcReq := &user.DeleteFriendReq{
		OwnerUserID:  req.OwnerUserID,
		FriendUserID: req.FriendUserID,
	}

	_, err = l.svcCtx.FriendClient.DeleteFriend(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteFriendResp{
	}

	return resp, nil
}
