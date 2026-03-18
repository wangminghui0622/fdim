package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyToAddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyToAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyToAddFriendLogic {
	return &ApplyToAddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyToAddFriendLogic) ApplyToAddFriend(req *user.ApplyToAddFriendReq) (*user.ApplyToAddFriendResp, error) {
	resp := &user.ApplyToAddFriendResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.FromUserID); err != nil {
		return nil, err
	}

	// ������֤
	if req.ToUserID == req.FromUserID {
		return nil, fmt.Errorf("can not add yourself")
	}

	// ����û��Ƿ����
	_, err := l.svcCtx.UserDB.FindWithError(l.ctx, []string{req.ToUserID, req.FromUserID})
	if err != nil {
		return nil, err
	}

	// ����Ƿ��Ѿ��Ǻ���
	in1, in2, err := l.svcCtx.FriendDB.CheckIn(l.ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, fmt.Errorf("already friends")
	}

	// Webhook BeforeAddFriend �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeAddFriendReq{
			CallbackCommand: webhook.CallbackBeforeAddFriendCommand,
			FromUserID:      req.FromUserID,
			ToUserID:        req.ToUserID,
			ReqMsg:          req.ReqMsg,
			Ex:              req.Ex,
		}
		cbResp := &webhook.CallbackBeforeAddFriendResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
	}

	// ��Ӻ�������
	if err := l.svcCtx.FriendDB.AddFriendRequest(l.ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}

	// 发送好友申请通知
	l.Infof("[ApplyToAddFriend] sending notification: from=%s to=%s, FriendNotification=%v",
		req.FromUserID, req.ToUserID, l.svcCtx.FriendNotification != nil)
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.FriendApplicationAddNotification(l.ctx, req.FromUserID, req.ToUserID)
		l.Infof("[ApplyToAddFriend] notification sent successfully: from=%s to=%s", req.FromUserID, req.ToUserID)
	} else {
		l.Errorf("[ApplyToAddFriend] FriendNotification is nil, cannot send notification")
	}

	// Webhook AfterAddFriend �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterAddFriendReq{
			CallbackCommand: webhook.CallbackAfterAddFriendCommand,
			FromUserID:      req.FromUserID,
			ToUserID:        req.ToUserID,
			ReqMsg:          req.ReqMsg,
		}
		cbResp := &webhook.CallbackAfterAddFriendResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
