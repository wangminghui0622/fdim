package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendListLogic) GetFriendList(req *types.GetFriendListReq) (resp *types.GetFriendListResp, err error) {
	rpcReq := &user.GetPaginationFriendsReq{
		UserID: req.UserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.FriendClient.GetPaginationFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var friendsInfo []types.FriendInfo
	for _, fi := range rpcResp.FriendsInfo {
		friendUserID := ""
		nickname := ""
		faceURL := ""
		var onlinePlatformIDs []int32
		
		if fi.FriendUser != nil {
			friendUserID = fi.FriendUser.UserID
			nickname = fi.FriendUser.Nickname
			faceURL = fi.FriendUser.FaceURL
			onlinePlatformIDs = fi.FriendUser.OnlinePlatformIDs
		}
		
		friendsInfo = append(friendsInfo, types.FriendInfo{
			OwnerUserID:       fi.OwnerUserID,
			FriendUserID:      friendUserID,
			Remark:            fi.Remark,
			CreateTime:        fi.CreateTime,
			AddSource:         fi.AddSource,
			Ex:                fi.Ex,
			Nickname:          nickname,
			FaceURL:           faceURL,
			OnlinePlatformIDs: onlinePlatformIDs,
		})
	}

	resp = &types.GetFriendListResp{
		FriendsInfo: friendsInfo,
		Total:       rpcResp.Total,
	}

	return resp, nil
}
