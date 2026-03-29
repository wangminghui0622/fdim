package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveBlackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveBlackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveBlackLogic {
	return &RemoveBlackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveBlackLogic) RemoveBlack(req *user.RemoveBlackReq) (*user.RemoveBlackResp, error) {
	resp := &user.RemoveBlackResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	if err := l.svcCtx.BlackDB.RemoveBlack(l.ctx, req.OwnerUserID, req.BlackUserID); err != nil {
		return nil, err
	}

	// ??????????????
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.BlackDeletedNotification(l.ctx, req.OwnerUserID, req.BlackUserID)
	}

	// Webhook AfterRemoveBlack ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterRemoveBlackReq{
			CallbackCommand: webhook.CallbackAfterRemoveBlackCommand,
			OwnerUserID:     req.OwnerUserID,
			BlackUserID:     req.BlackUserID,
		}
		cbResp := &webhook.CallbackAfterRemoveBlackResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
