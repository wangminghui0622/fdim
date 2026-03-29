package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/pkg/notification"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetConversationHasReadSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetConversationHasReadSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationHasReadSeqLogic {
	return &SetConversationHasReadSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetConversationHasReadSeqLogic) SetConversationHasReadSeq(req *msg.SetConversationHasReadSeqReq) (*msg.SetConversationHasReadSeqResp, error) {
	if req.ConversationID == "" || req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID and userID are required")
	}

	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	if req.HasReadSeq == 0 {
		return &msg.SetConversationHasReadSeqResp{}, nil
	}

	if err := l.svcCtx.MsgCache.SetHasReadSeqs(l.ctx, req.ConversationID, map[string]int64{
		req.UserID: req.HasReadSeq,
	}); err != nil {
		l.Errorw("SetHasReadSeqs failed", logx.Field("conversationID", req.ConversationID), logx.Field("userID", req.UserID), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to set has read seq")
	}

	// 发送会话未读数变更通知给自己（与官方一致，除非明确禁用?
	if !req.NoNotification && l.svcCtx.NotificationSender != nil {
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, req.ConversationID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("error", err))
		} else {
			unreadCount := int64(0)
			if maxSeq > req.HasReadSeq {
				unreadCount = maxSeq - req.HasReadSeq
			}
			// 使用 ConversationNotificationSender 包装?
			convSender := &notification.ConversationNotificationSender{NotificationSender: l.svcCtx.NotificationSender}
			go convSender.ConversationUnreadChangeNotification(
				l.ctx, req.UserID, req.ConversationID, unreadCount, req.HasReadSeq)
		}
	}

	return &msg.SetConversationHasReadSeqResp{}, nil
}
