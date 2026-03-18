package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/mcontext"
	"fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// ConversationNotificationSender 会话通知发送器
type ConversationNotificationSender struct {
	*NotificationSender
}

// NewConversationNotificationSender 创建会话通知发送器
func NewConversationNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *ConversationNotificationSender {
	return &ConversationNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// ConversationChangeNotification 发送会话变更通知
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string, conversationIDs []string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = userID
	}
	tips := &sdkws.ConversationUpdateTips{
		UserID:             userID,
		ConversationIDList: conversationIDs,
	}
	c.Notification(ctx, opUserID, userID, constant.ConversationChangeNotification, tips)
}

// ConversationUnreadChangeNotification 发送会话未读数变更通知
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(ctx context.Context, userID string, conversationID string, unreadCount int64, hasReadSeq int64) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = userID
	}
	tips := &sdkws.ConversationHasReadTips{
		UserID:         userID,
		ConversationID: conversationID,
		HasReadSeq:     hasReadSeq,
	}
	c.Notification(ctx, opUserID, userID, constant.ConversationUnreadNotification, tips)
}
