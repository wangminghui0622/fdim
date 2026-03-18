package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type UpdateFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendsLogic {
	return &UpdateFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFriendsLogic) UpdateFriends(req *types.UpdateFriendsReq) (resp *types.UpdateFriendsResp, err error) {
	rpcReq := &user.UpdateFriendsReq{
		OwnerUserID:   req.OwnerUserID,
		FriendUserIDs: req.FriendUserIDs,
	}

	_, err = l.svcCtx.FriendClient.UpdateFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateFriendsResp{
	}

	return resp, nil
}
