package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/model"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RespondFriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRespondFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RespondFriendApplyLogic {
	return &RespondFriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RespondFriendApplyLogic) RespondFriendApply(req *user.RespondFriendApplyReq) (*user.RespondFriendApplyResp, error) {
	resp := &user.RespondFriendApplyResp{}

	// Ȩ֤
	if err := authverify.CheckAccess(l.ctx, req.ToUserID); err != nil {
		return nil, err
	}

	friendRequest := model.FriendRequest{
		FromUserID:    req.FromUserID,
		ToUserID:      req.ToUserID,
		HandleMsg:     req.HandleMsg,
		HandleResult:  req.HandleResult,
		HandlerUserID: req.ToUserID, // ǽշ
		HandleTime:    time.Now(),
	}

	// ֤
	if req.HandleResult != constant.FriendRequestAgree && req.HandleResult != constant.FriendRequestRefuse {
		return nil, fmt.Errorf("invalid handle result: %d", req.HandleResult)
	}

	if req.HandleResult == constant.FriendRequestAgree {
		// Webhook BeforeAddFriendAgree ص
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackBeforeAddFriendAgreeReq{
				CallbackCommand: webhook.CallbackBeforeAddFriendAgreeCommand,
				FromUserID:      req.FromUserID,
				ToUserID:        req.ToUserID,
				HandleMsg:       req.HandleMsg,
				HandleResult:    req.HandleResult,
			}
			cbResp := &webhook.CallbackBeforeAddFriendAgreeResp{}
			if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
				if err != webhook.ErrCallbackContinue {
					return nil, err
				}
				// ErrCallbackContinue ʾִ
			}
		}

		if err := l.svcCtx.FriendDB.AgreeFriendRequest(l.ctx, &friendRequest); err != nil {
			return nil, err
		}

		if l.svcCtx.FriendNotification != nil {
			l.svcCtx.FriendNotification.FriendApplicationAgreedNotification(l.ctx, req.FromUserID, req.ToUserID, req.HandleMsg)
		}

		// Webhook AfterAddFriendAgree ص
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackAfterAddFriendAgreeReq{
				CallbackCommand: webhook.CallbackAfterAddFriendAgreeCommand,
				FromUserID:      req.FromUserID,
				ToUserID:        req.ToUserID,
				HandleMsg:       req.HandleMsg,
			}
			cbResp := &webhook.CallbackAfterAddFriendAgreeResp{}
			l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
		}
	} else {
		if err := l.svcCtx.FriendDB.RefuseFriendRequest(l.ctx, &friendRequest); err != nil {
			return nil, err
		}
		if l.svcCtx.FriendNotification != nil {
			l.svcCtx.FriendNotification.FriendApplicationRejectedNotification(l.ctx, req.FromUserID, req.ToUserID, req.HandleMsg)
		}
	}

	return resp, nil
}
