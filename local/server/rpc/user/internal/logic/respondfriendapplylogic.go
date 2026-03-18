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

	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.ToUserID); err != nil {
		return nil, err
	}

	friendRequest := model.FriendRequest{
		FromUserID:    req.FromUserID,
		ToUserID:      req.ToUserID,
		HandleMsg:     req.HandleMsg,
		HandleResult:  req.HandleResult,
		HandlerUserID: req.ToUserID, // 处理者是接收方
		HandleTime:    time.Now(),
	}

	// 参数验证
	if req.HandleResult != constant.FriendRequestAgree && req.HandleResult != constant.FriendRequestRefuse {
		return nil, fmt.Errorf("invalid handle result: %d", req.HandleResult)
	}

	if req.HandleResult == constant.FriendRequestAgree {
		// Webhook BeforeAddFriendAgree 回调
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
				// ErrCallbackContinue 表示继续执行
			}
		}

		if err := l.svcCtx.FriendDB.AgreeFriendRequest(l.ctx, &friendRequest); err != nil {
			return nil, err
		}

		// 发送好友申请同意通知
		// 注意：通知消息使用 SingleChatType，会自动触发会话创建（在 msg.SendMsg -> ensureConversation 中）
		// 这是官方的实现方式，不需要在这里手动创建会话
		if l.svcCtx.FriendNotification != nil {
			l.svcCtx.FriendNotification.FriendApplicationAgreedNotification(l.ctx, req.FromUserID, req.ToUserID, req.HandleMsg)
		}

		// Webhook AfterAddFriendAgree 回调
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
		// 拒绝
		if err := l.svcCtx.FriendDB.RefuseFriendRequest(l.ctx, &friendRequest); err != nil {
			return nil, err
		}
		// 发送好友申请拒绝通知
		if l.svcCtx.FriendNotification != nil {
			l.svcCtx.FriendNotification.FriendApplicationRejectedNotification(l.ctx, req.FromUserID, req.ToUserID, req.HandleMsg)
		}
		// 注意：拒绝后不需要通知，也不需要 webhook 回调（参考 open-im-server 的实现）
	}

	return resp, nil
}
