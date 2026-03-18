package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDesignatedFriendsApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDesignatedFriendsApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesignatedFriendsApplyLogic {
	return &GetDesignatedFriendsApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDesignatedFriendsApplyLogic) GetDesignatedFriendsApply(req *user.GetDesignatedFriendsApplyReq) (*user.GetDesignatedFriendsApplyResp, error) {
	resp := &user.GetDesignatedFriendsApplyResp{}

	// Ȩ����֤����Ҫ�������ϵ�е�һ�������Ա
	if err := authverify.CheckAccessIn(l.ctx, req.FromUserID, req.ToUserID); err != nil {
		return nil, err
	}

	// ����˫���������
	friendRequests, err := l.svcCtx.FriendDB.FindBothFriendRequests(l.ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}

	// ��ȡ������ص��û�ID�����ͷ��ͽ��շ���
	userIDs := make(map[string]bool)
	for _, req := range friendRequests {
		userIDs[req.FromUserID] = true
		userIDs[req.ToUserID] = true
	}

	// ת��Ϊ��Ƭ
	userIDList := make([]string, 0, len(userIDs))
	for userID := range userIDs {
		userIDList = append(userIDList, userID)
	}

	// ��ȡ�û���Ϣ
	users, err := l.svcCtx.UserDB.Find(l.ctx, userIDList)
	if err != nil {
		return nil, err
	}

	// �����û�ӳ��
	userMap := make(map[string]*model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// 转换为 Protocol Buffer 格式并填充用户信息
	result := make([]*sdkws.FriendRequest, 0, len(friendRequests))
	for _, req := range friendRequests {
		friendRequest := convert.ModelFriendRequestDB2Pb(req)

		// ��䷢�ͷ��û���Ϣ
		if fromUser, ok := userMap[req.FromUserID]; ok {
			friendRequest.FromNickname = fromUser.Nickname
			friendRequest.FromFaceURL = fromUser.FaceURL
		}

		// �����շ��û���Ϣ
		if toUser, ok := userMap[req.ToUserID]; ok {
			friendRequest.ToNickname = toUser.Nickname
			friendRequest.ToFaceURL = toUser.FaceURL
		}

		result = append(result, friendRequest)
	}

	resp.FriendRequests = result
	return resp, nil
}
