package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetSpecifiedFriendsInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSpecifiedFriendsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedFriendsInfoLogic {
	return &GetSpecifiedFriendsInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecifiedFriendsInfoLogic) GetSpecifiedFriendsInfo(req *types.GetSpecifiedFriendsInfoReq) (resp *types.GetSpecifiedFriendsInfoResp, err error) {
	rpcReq := &user.GetSpecifiedFriendsInfoReq{
		OwnerUserID: req.OwnerUserID,
		UserIDList:  req.UserIDList,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetSpecifiedFriendsInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var friendsInfo []types.FriendInfo
	for _, info := range rpcResp.Infos {
		if info.FriendInfo == nil {
			continue
		}
		friendUserID := ""
		if info.FriendInfo.FriendUser != nil {
			friendUserID = info.FriendInfo.FriendUser.UserID
		}
		friendsInfo = append(friendsInfo, types.FriendInfo{
			OwnerUserID:  info.FriendInfo.OwnerUserID,
			FriendUserID: friendUserID,
			Remark:       info.FriendInfo.Remark,
			CreateTime:   info.FriendInfo.CreateTime,
			AddSource:    info.FriendInfo.AddSource,
			Ex:           info.FriendInfo.Ex,
		})
	}

	resp = &types.GetSpecifiedFriendsInfoResp{
		FriendsInfo: friendsInfo,
	}

	return resp, nil
}
