package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ImportFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportFriendsLogic {
	return &ImportFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportFriendsLogic) ImportFriends(req *types.ImportFriendsReq) (resp *types.ImportFriendsResp, err error) {
	rpcReq := &user.ImportFriendReq{
		OwnerUserID:   req.OwnerUserID,
		FriendUserIDs: req.FriendUserIDs,
	}

	_, err = l.svcCtx.FriendClient.ImportFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ImportFriendsResp{
	}

	return resp, nil
}
