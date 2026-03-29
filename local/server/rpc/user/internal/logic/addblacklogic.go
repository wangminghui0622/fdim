package logic

import (
	"context"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddBlackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBlackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBlackLogic {
	return &AddBlackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBlackLogic) AddBlack(req *user.AddBlackReq) (*user.AddBlackResp, error) {
	resp := &user.AddBlackResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// Webhook BeforeAddBlack ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeAddBlackReq{
			CallbackCommand: webhook.CallbackBeforeAddBlackCommand,
			OwnerUserID:     req.OwnerUserID,
			BlackUserID:     req.BlackUserID,
		}
		cbResp := &webhook.CallbackBeforeAddBlackResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ??????????
		}
	}

	// ??????????ID
	opUserID := mcontext.GetOpUserID(l.ctx)
	if opUserID == "" {
		opUserID = req.OwnerUserID
	}

	black := &model.Black{
		OwnerUserID:    req.OwnerUserID,
		BlockUserID:    req.BlackUserID,
		CreateTime:     time.Now(),
		AddSource:      0, // AddBlackReq ????? AddSource ??????????
		OperatorUserID: opUserID,
		Ex:             req.Ex,
	}

	if err := l.svcCtx.BlackDB.AddBlack(l.ctx, black); err != nil {
		return nil, err
	}

	// ??????????????
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.BlackAddedNotification(l.ctx, req.OwnerUserID, req.BlackUserID)
	}

	// Webhook AfterAddBlack ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterAddBlackReq{
			CallbackCommand: webhook.CallbackAfterAddBlackCommand,
			OwnerUserID:     req.OwnerUserID,
			BlackUserID:     req.BlackUserID,
		}
		cbResp := &webhook.CallbackAfterAddBlackResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
