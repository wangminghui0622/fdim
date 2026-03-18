package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/constant"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarkMsgsAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkMsgsAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkMsgsAsReadLogic {
	return &MarkMsgsAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkMsgsAsReadLogic) MarkMsgsAsRead(req *msg.MarkMsgsAsReadReq) (*msg.MarkMsgsAsReadResp, error) {
	if req.ConversationID == "" || req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID and userID are required")
	}
	if len(req.Seqs) == 0 {
		return &msg.MarkMsgsAsReadResp{}, nil
	}

	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	// 取最大 seq 作为该用户的已读 seq
	var maxSeq int64
	for _, s := range req.Seqs {
		if s > maxSeq {
			maxSeq = s
		}
	}
	if maxSeq == 0 {
		return &msg.MarkMsgsAsReadResp{}, nil
	}

	if err := l.svcCtx.MsgCache.SetHasReadSeqs(l.ctx, req.ConversationID, map[string]int64{
		req.UserID: maxSeq,
	}); err != nil {
		l.Errorw("SetHasReadSeqs failed", logx.Field("conversationID", req.ConversationID), logx.Field("userID", req.UserID), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to set has read seqs")
	}

	// 发送 HasReadReceipt 通知（与官方一致）
	if l.svcCtx.NotificationSender != nil {
		go l.sendMarkAsReadNotification(req.UserID, req.ConversationID, req.Seqs, maxSeq)
	}

	return &msg.MarkMsgsAsReadResp{}, nil
}

// sendMarkAsReadNotification 发送已读回执通知（完全参考官方 as_read.go 实现）
func (l *MarkMsgsAsReadLogic) sendMarkAsReadNotification(userID, conversationID string, seqs []int64, hasReadSeq int64) {
	// 获取会话信息以确定接收方（与官方一致）
	convResp, err := l.svcCtx.ConversationClient.GetConversation(l.ctx, &conversation.GetConversationReq{
		ConversationID: conversationID,
		OwnerUserID:    userID,
	})
	if err != nil {
		l.Errorw("get conversation failed", logx.Field("conversationID", conversationID), logx.Field("error", err))
		return
	}

	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: userID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}

	// 确定接收方（与官方 conversationAndGetRecvID 逻辑一致）
	var recvID string
	if convResp.Conversation.ConversationType == constant.SingleChatType ||
		convResp.Conversation.ConversationType == constant.NotificationChatType {
		// 单聊：发给对方
		if userID == convResp.Conversation.OwnerUserID {
			recvID = convResp.Conversation.UserID
		} else {
			recvID = convResp.Conversation.OwnerUserID
		}
	} else if convResp.Conversation.ConversationType == constant.ReadGroupChatType {
		// 群聊：发给群
		recvID = convResp.Conversation.GroupID
	}

	// 使用 NotificationSender 发送通知（与官方一致）
	l.svcCtx.NotificationSender.NotificationWithSessionType(
		l.ctx,
		userID,
		recvID,
		constant.HasReadReceipt,
		convResp.Conversation.ConversationType,
		tips,
	)
}
