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

	// ������֤
	if len(req.FriendUserIDs) == 0 {
		return nil, fmt.Errorf("friendUserIDs is empty")
	}

	// ����ظ��� FriendUserIDs
	seen := make(map[string]bool)
	for _, userID := range req.FriendUserIDs {
		if seen[userID] {
			return nil, fmt.Errorf("friend userID repeated: %s", userID)
		}
		seen[userID] = true
	}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ���Һ���
	friends, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// ��ȡ���ѵ��û���Ϣ
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// �����û�ӳ��
	userMap := make(map[string]*model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// ת��Ϊ Protocol Buffer ��ʽ������û���Ϣ
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

		// �����ѵ��û���Ϣ
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
