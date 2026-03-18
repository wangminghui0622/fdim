package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetDesignatedFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDesignatedFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesignatedFriendsLogic {
	return &GetDesignatedFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDesignatedFriendsLogic) GetDesignatedFriends(req *types.GetDesignatedFriendsReq) (resp *types.GetDesignatedFriendsResp, err error) {
	rpcReq := &user.GetDesignatedFriendsReq{
		OwnerUserID:   req.OwnerUserID,
		FriendUserIDs: req.FriendUserIDs,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetDesignatedFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var friendsInfo []types.FriendInfo
	for _, fi := range rpcResp.FriendsInfo {
		friendUserID := ""
		if fi.FriendUser != nil {
			friendUserID = fi.FriendUser.UserID
		}
		friendsInfo = append(friendsInfo, types.FriendInfo{
			OwnerUserID:  fi.OwnerUserID,
			FriendUserID: friendUserID,
			Remark:       fi.Remark,
			CreateTime:   fi.CreateTime,
			AddSource:    fi.AddSource,
			Ex:           fi.Ex,
		})
	}

	resp = &types.GetDesignatedFriendsResp{
		FriendsInfo: friendsInfo,
	}

	return resp, nil
}
