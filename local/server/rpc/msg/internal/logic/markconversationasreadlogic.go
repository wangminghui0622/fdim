package logic

import (
	"context"
	"errors"

	"fdim/pkg/notification"
	"fdim/protocol/constant"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarkConversationAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkConversationAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkConversationAsReadLogic {
	return &MarkConversationAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkConversationAsReadLogic) MarkConversationAsRead(req *msg.MarkConversationAsReadReq) (*msg.MarkConversationAsReadResp, error) {
	l.Infow("MarkConversationAsRead called",
		logx.Field("userID", req.UserID),
		logx.Field("conversationID", req.ConversationID),
		logx.Field("hasReadSeq", req.HasReadSeq),
		logx.Field("seqs", req.Seqs))

	// 获取当前已读序列号
	hasReadSeq, err := l.svcCtx.MsgDatabase.GetHasReadSeq(l.ctx, req.UserID, req.ConversationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// 检查 ConversationClient 是否初始化
	if l.svcCtx.ConversationClient == nil {
		l.Errorw("ConversationClient is nil, cannot get conversation info")
		// 如果没有 ConversationClient，直接更新已读序列号
		if req.HasReadSeq > hasReadSeq {
			err = l.svcCtx.MsgDatabase.SetHasReadSeq(l.ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				return nil, err
			}
		}
		
		return &msg.MarkConversationAsReadResp{}, nil
	}
	
	// 【修复】获取会话信息，如果会话不存在则只更新已读序列号（不返回错误）
	conversationResp, err := l.svcCtx.ConversationClient.GetConversation(l.ctx, &conversation.GetConversationReq{
		ConversationID: req.ConversationID,
		OwnerUserID:    req.UserID,
	})
	if err != nil {
		l.Infow("GetConversation failed, conversation may not exist yet, only update hasReadSeq",
			logx.Field("error", err),
			logx.Field("conversationID", req.ConversationID),
			logx.Field("userID", req.UserID))
		
		// 会话不存在时，只更新已读序列号，不发送通知
		if req.HasReadSeq > hasReadSeq {
			err = l.svcCtx.MsgDatabase.SetHasReadSeq(l.ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				return nil, err
			}
		}
		
		return &msg.MarkConversationAsReadResp{}, nil
	}
	
	var seqs []int64

	l.Infow("Current hasReadSeq from DB",
		logx.Field("hasReadSeq", hasReadSeq),
		logx.Field("req.HasReadSeq", req.HasReadSeq),
		logx.Field("conversationType", conversationResp.Conversation.ConversationType))
	
	if conversationResp.Conversation.ConversationType == constant.SingleChatType {
		// 生成需要标记为已读的序列号列表
		for i := hasReadSeq + 1; i <= req.HasReadSeq; i++ {
			seqs = append(seqs, i)
		}
		// 避免客户端乱序调用导致遗漏
		for _, val := range req.Seqs {
			if !contains(val, seqs) {
				seqs = append(seqs, val)
			}
		}
		if len(seqs) > 0 {
			l.Infow("Marking messages as read", logx.Field("seqs", seqs), logx.Field("conversationID", req.ConversationID))
			if err = l.svcCtx.MsgDatabase.MarkSingleChatMsgsAsRead(l.ctx, req.UserID, req.ConversationID, seqs); err != nil {
				l.Errorw("MarkSingleChatMsgsAsRead failed", logx.Field("error", err))
				return nil, err
			}
		}
		if req.HasReadSeq > hasReadSeq {
			err = l.svcCtx.MsgDatabase.SetHasReadSeq(l.ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				l.Errorw("SetHasReadSeq failed", logx.Field("error", err))
				return nil, err
			}
			hasReadSeq = req.HasReadSeq
			l.Infow("Updated hasReadSeq", logx.Field("newHasReadSeq", hasReadSeq))
		}
		
		// 【修复问题1】发送已读回执给消息发送方（对方用户）
		// 这样发送方可以看到消息已被对方读取
		recvID := l.conversationAndGetRecvID(conversationResp.Conversation, req.UserID)
		l.Infow("Sending HasReadReceipt notification to sender",
			logx.Field("fromUserID", req.UserID),
			logx.Field("toUserID", recvID),
			logx.Field("seqs", seqs),
			logx.Field("hasReadSeq", hasReadSeq))
		l.sendMarkAsReadNotification(l.ctx, req.ConversationID, conversationResp.Conversation.ConversationType, req.UserID, recvID, seqs, hasReadSeq)
	} else if conversationResp.Conversation.ConversationType == constant.ReadGroupChatType ||
		conversationResp.Conversation.ConversationType == constant.NotificationChatType {
		if req.HasReadSeq > hasReadSeq {
			err = l.svcCtx.MsgDatabase.SetHasReadSeq(l.ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				return nil, err
			}
			hasReadSeq = req.HasReadSeq
		}
		l.sendMarkAsReadNotification(l.ctx, req.ConversationID, constant.SingleChatType, req.UserID,
			req.UserID, seqs, hasReadSeq)
	}

	// 【修复问题2】发送会话未读数变更通知给自己（与官方一致）
	// 这样自己的会话列表会更新未读数为0
	if l.svcCtx.NotificationSender != nil {
		maxSeq, err := l.svcCtx.MsgDatabase.GetMaxSeq(l.ctx, req.ConversationID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("error", err))
		} else {
			unreadCount := int64(0)
			if maxSeq > hasReadSeq {
				unreadCount = maxSeq - hasReadSeq
			}
			l.Infow("Sending ConversationUnreadChange notification to self",
				logx.Field("userID", req.UserID),
				logx.Field("conversationID", req.ConversationID),
				logx.Field("unreadCount", unreadCount),
				logx.Field("hasReadSeq", hasReadSeq))
			// 使用 ConversationNotificationSender 包装器
			convSender := &notification.ConversationNotificationSender{NotificationSender: l.svcCtx.NotificationSender}
			go convSender.ConversationUnreadChangeNotification(
				l.ctx, req.UserID, req.ConversationID, unreadCount, hasReadSeq)
		}
	}

	l.Infow("MarkConversationAsRead completed successfully",
		logx.Field("userID", req.UserID),
		logx.Field("conversationID", req.ConversationID),
		logx.Field("finalHasReadSeq", hasReadSeq))

	return &msg.MarkConversationAsReadResp{}, nil
}

func (l *MarkConversationAsReadLogic) conversationAndGetRecvID(conversation *conversation.Conversation, userID string) string {
	if conversation.ConversationType == constant.SingleChatType ||
		conversation.ConversationType == constant.NotificationChatType {
		if userID == conversation.OwnerUserID {
			return conversation.UserID
		} else {
			return conversation.OwnerUserID
		}
	} else if conversation.ConversationType == constant.ReadGroupChatType {
		return conversation.GroupID
	}
	return ""
}

func (l *MarkConversationAsReadLogic) sendMarkAsReadNotification(ctx context.Context, conversationID string, sessionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	l.svcCtx.NotificationSender.NotificationWithSessionType(ctx, sendID, recvID, constant.HasReadReceipt, sessionType, tips)
}

func contains(val int64, seqs []int64) bool {
	for _, s := range seqs {
		if s == val {
			return true
		}
	}
	return false
}
