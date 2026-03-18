package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ImportFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportFriendsLogic {
	return &ImportFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ImportFriendsLogic) ImportFriends(req *user.ImportFriendReq) (*user.ImportFriendResp, error) {
	resp := &user.ImportFriendResp{}

	// Ȩ����֤����Ҫ����ԱȨ��
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// ������֤
	if len(req.FriendUserIDs) == 0 {
		return nil, fmt.Errorf("friendUserIDs is empty")
	}

	// ����Ƿ�����Լ�
	for _, userID := range req.FriendUserIDs {
		if userID == req.OwnerUserID {
			return nil, fmt.Errorf("can not add yourself")
		}
	}

	// ��� friendUserIDs �Ƿ��ظ�
	userIDMap := make(map[string]bool)
	for _, userID := range req.FriendUserIDs {
		if userIDMap[userID] {
			return nil, fmt.Errorf("friend userID repeated: %s", userID)
		}
		userIDMap[userID] = true
	}

	// ����û��Ƿ����
	allUserIDs := append([]string{req.OwnerUserID}, req.FriendUserIDs...)
	_, err := l.svcCtx.UserDB.FindWithError(l.ctx, allUserIDs)
	if err != nil {
		return nil, err
	}

	// ����Ƿ��Ѿ��Ǻ��ѣ�����Ѿ��Ǻ��ѣ�������
	existingFriends, err := l.svcCtx.FriendDB.FindFriends(l.ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// ���˵��Ѿ��Ǻ��ѵ��û�
	existingFriendMap := make(map[string]bool)
	for _, friend := range existingFriends {
		existingFriendMap[friend.FriendUserID] = true
	}

	newFriendUserIDs := make([]string, 0, len(req.FriendUserIDs))
	for _, userID := range req.FriendUserIDs {
		if !existingFriendMap[userID] {
			newFriendUserIDs = append(newFriendUserIDs, userID)
		}
	}

	if len(newFriendUserIDs) == 0 {
		// �����û����Ѿ��Ǻ��ѣ�ֱ�ӷ���
		return resp, nil
	}

	// Webhook BeforeImportFriends �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeImportFriendsReq{
			CallbackCommand: webhook.CallbackBeforeImportFriendsCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserIDs:   newFriendUserIDs,
		}
		cbResp := &webhook.CallbackBeforeImportFriendsResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// ��� webhook �������޸ĺ�ĺ����б��ʹ����
		if len(cbResp.FriendUserIDs) > 0 {
			newFriendUserIDs = cbResp.FriendUserIDs
		}
	}

	// ������ѹ�ϵ
	if err := l.svcCtx.FriendDB.BecomeFriends(l.ctx, req.OwnerUserID, newFriendUserIDs, constant.BecomeFriendByImport); err != nil {
		return nil, err
	}

	// ���ͺ�������ͬ��֪ͨ����������൱���Զ�ͬ�⣩
	if l.svcCtx.FriendNotification != nil {
		for _, friendUserID := range newFriendUserIDs {
			l.svcCtx.FriendNotification.FriendAddedNotification(l.ctx, req.OwnerUserID, friendUserID)
			l.svcCtx.FriendNotification.FriendAddedNotification(l.ctx, friendUserID, req.OwnerUserID)
		}
	}

	// Webhook AfterImportFriends �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterImportFriendsReq{
			CallbackCommand: webhook.CallbackAfterImportFriendsCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserIDs:   newFriendUserIDs,
		}
		cbResp := &webhook.CallbackAfterImportFriendsResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
