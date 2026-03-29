package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSpecifiedFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSpecifiedFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedFriendsLogic {
	return &GetSpecifiedFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSpecifiedFriendsLogic) GetSpecifiedFriends(req *user.GetSpecifiedFriendsInfoReq) (*user.GetSpecifiedFriendsInfoResp, error) {
	resp := &user.GetSpecifiedFriendsInfoResp{}

	// ???????
	if len(req.UserIDList) == 0 {
		return nil, fmt.Errorf("userIDList is empty")
	}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ???????
	friends, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}

	// ??????????????
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.UserIDList)
	if err != nil {
		return nil, err
	}

	// ??????????
	userMap := make(map[string]*model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// ???? Protocol Buffer ??????????????
	result := make([]*user.GetSpecifiedFriendsInfoInfo, 0, len(req.UserIDList))
	for _, userID := range req.UserIDList {
		info := &user.GetSpecifiedFriendsInfoInfo{}

		// ?????????
		if user, ok := userMap[userID]; ok {
			info.UserInfo = &sdkws.UserInfo{
				UserID:   user.UserID,
				Nickname: user.Nickname,
				FaceURL:  user.FaceURL,
				Ex:       user.Ex,
			}
		}

		// Ϣڣ
		for _, friend := range friends {
			if friend.FriendUserID == userID {
				relationFriend := convert.ModelFriendDB2Pb(friend, nil)
				info.FriendInfo = &sdkws.FriendInfo{
					OwnerUserID:    relationFriend.OwnerUserID,
					Remark:         relationFriend.Remark,
					CreateTime:     relationFriend.CreateTime,
					AddSource:      relationFriend.AddSource,
					OperatorUserID: relationFriend.OperatorUserID,
					Ex:             relationFriend.Ex,
					IsPinned:       relationFriend.IsPinned,
				}
				if info.UserInfo != nil {
					info.FriendInfo.FriendUser = info.UserInfo
				}
				break
			}
		}

		result = append(result, info)
	}

	resp.Infos = result

	return resp, nil
}
