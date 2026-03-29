package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDesignatedFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDesignatedFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesignatedFriendsLogic {
	return &GetDesignatedFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDesignatedFriendsLogic) GetDesignatedFriends(req *user.GetDesignatedFriendsReq) (*user.GetDesignatedFriendsResp, error) {
	resp := &user.GetDesignatedFriendsResp{}

	// ???????
	if len(req.FriendUserIDs) == 0 {
		return nil, fmt.Errorf("friendUserIDs is empty")
	}

	// ???????? FriendUserIDs
	seen := make(map[string]bool)
	for _, userID := range req.FriendUserIDs {
		if seen[userID] {
			return nil, fmt.Errorf("friend userID repeated: %s", userID)
		}
		seen[userID] = true
	}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ???????
	friends, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// ??????????????
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// ??????????
	userMap := make(map[string]*model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// ???? Protocol Buffer ??????????????
	result := make([]*sdkws.FriendInfo, 0, len(friends))
	for _, friend := range friends {
		friendInfo := &sdkws.FriendInfo{
			OwnerUserID:    friend.OwnerUserID,
			Remark:         friend.Remark,
			CreateTime:     friend.CreateTime.UnixMilli(),
			AddSource:      friend.AddSource,
			OperatorUserID: friend.OperatorUserID,
			Ex:             friend.Ex,
			IsPinned:       friend.IsPinned,
		}

		// ?????????????
		if user, ok := userMap[friend.FriendUserID]; ok {
			friendInfo.FriendUser = &sdkws.UserInfo{
				UserID:   user.UserID,
				Nickname: user.Nickname,
				FaceURL:  user.FaceURL,
				Ex:       user.Ex,
			}
		}

		result = append(result, friendInfo)
	}

	resp.FriendsInfo = result
	return resp, nil
}
