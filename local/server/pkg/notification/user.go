package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/mcontext"
	"fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// UserNotificationSender û֪ͨ
type UserNotificationSender struct {
	*NotificationSender
}

// NewUserNotificationSender û֪ͨ
func NewUserNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *UserNotificationSender {
	return &UserNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// UserInfoUpdatedNotification ûϢ֪ͨ
func (u *UserNotificationSender) UserInfoUpdatedNotification(ctx context.Context, changedUserID string, needNotifiedUserID string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = changedUserID
	}
	tips := &sdkws.UserInfoUpdatedTips{
		UserID: changedUserID,
	}
	u.Notification(ctx, opUserID, needNotifiedUserID, constant.UserInfoUpdatedNotification, tips)
}

// UserStatusChangeNotification û״̬֪ͨ
func (u *UserNotificationSender) UserStatusChangeNotification(ctx context.Context, tips *sdkws.UserStatusChangeTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserStatusChangeNotification, tips)
}

// UserCommandAddNotification û֪ͨ
func (u *UserNotificationSender) UserCommandAddNotification(ctx context.Context, tips *sdkws.UserCommandAddTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandAddNotification, tips)
}

// UserCommandDeleteNotification ûɾ֪ͨ
func (u *UserNotificationSender) UserCommandDeleteNotification(ctx context.Context, tips *sdkws.UserCommandDeleteTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandDeleteNotification, tips)
}

// UserCommandUpdateNotification û֪ͨ
func (u *UserNotificationSender) UserCommandUpdateNotification(ctx context.Context, tips *sdkws.UserCommandUpdateTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandUpdateNotification, tips)
}
