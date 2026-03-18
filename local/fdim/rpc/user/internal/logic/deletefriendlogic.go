package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFriendLogic) DeleteFriend(req *user.DeleteFriendReq) (*user.DeleteFriendResp, error) {
	resp := &user.DeleteFriendResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ����Ƿ�Ϊ����
	_, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}

	// ɾ�����ѹ�ϵ
	if err := l.svcCtx.FriendDB.Delete(l.ctx, req.OwnerUserID, []string{req.FriendUserID}); err != nil {
		return nil, err
	}

	// ���ͺ���ɾ��֪ͨ
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.FriendDeletedNotification(l.ctx, req.OwnerUserID, req.FriendUserID)
	}

	// Webhook AfterDeleteFriend �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterDeleteFriendReq{
			CallbackCommand: webhook.CallbackAfterDeleteFriendCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserID:    req.FriendUserID,
		}
		cbResp := &webhook.CallbackAfterDeleteFriendResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
