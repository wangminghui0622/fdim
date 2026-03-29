package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaginationFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaginationFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaginationFriendsLogic {
	return &GetPaginationFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaginationFriendsLogic) GetPaginationFriends(req *user.GetPaginationFriendsReq) (*user.GetPaginationFriendsResp, error) {
	resp := &user.GetPaginationFriendsResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ??????????
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ????????????
	total, friends, err := l.svcCtx.FriendDB.PageOwnerFriends(l.ctx, req.UserID, offset, limit)
	if err != nil {
		return nil, err
	}

	// תΪ Protocol Buffer ʽ
	resp.Total = int32(total)
	// ת Friend Ϊ sdkws.FriendInfo
	relationFriends := convert.ModelFriendsDB2Pb(friends)
	resp.FriendsInfo = make([]*sdkws.FriendInfo, 0, len(relationFriends))
	for _, rf := range relationFriends {
		sdkwsFriend := &sdkws.FriendInfo{
			OwnerUserID:    rf.OwnerUserID,
			Remark:         rf.Remark,
			CreateTime:     rf.CreateTime,
			AddSource:      rf.AddSource,
			OperatorUserID: rf.OperatorUserID,
			Ex:             rf.Ex,
			IsPinned:       rf.IsPinned,
		}
		resp.FriendsInfo = append(resp.FriendsInfo, sdkwsFriend)
	}

	// ??????????????
	if len(friends) > 0 {
		friendUserIDs := make([]string, 0, len(friends))
		for _, f := range friends {
			friendUserIDs = append(friendUserIDs, f.FriendUserID)
		}

		users, err := l.svcCtx.UserDB.Find(l.ctx, friendUserIDs)
		if err != nil {
			return nil, err
		}

		// ??????????
		userMap := make(map[string]*model.User)
		for _, u := range users {
			userMap[u.UserID] = u
		}

		// ѯѵ״̬
		onlineStatusMap := make(map[string][]int32)
		for _, friendUserID := range friendUserIDs {
			platformIDs, err := l.svcCtx.UserCache.GetUserOnline(l.ctx, friendUserID)
			if err != nil {
				// ¼󵫲ж
				l.Logger.Errorf("Failed to get online status for user %s: %v", friendUserID, err)
				continue
			}
			if len(platformIDs) > 0 {
				onlineStatusMap[friendUserID] = platformIDs
			}
		}

		// ?????????????
		// ???? convert.FriendsDB2Pb ????????????????? friends ?????????????????????????
		for i, friendInfo := range resp.FriendsInfo {
			if i < len(friends) {
				friend := friends[i]
				if user, ok := userMap[friend.FriendUserID]; ok {
					userInfo := &sdkws.UserInfo{
						UserID:   user.UserID,
						Nickname: user.Nickname,
						FaceURL:  user.FaceURL,
						Ex:       user.Ex,
					}
					
					// ״̬Ϣ
					if platformIDs, exists := onlineStatusMap[friend.FriendUserID]; exists && len(platformIDs) > 0 {
						// ûߣƽ̨б
						userInfo.OnlinePlatformIDs = platformIDs
					}
					
					friendInfo.FriendUser = userInfo
				}
			}
		}
	}

	return resp, nil
}
