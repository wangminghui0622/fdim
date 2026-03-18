package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/mcontext"
	"fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// UserNotificationSender 用户通知发送器
type UserNotificationSender struct {
	*NotificationSender
}

// NewUserNotificationSender 创建用户通知发送器
func NewUserNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *UserNotificationSender {
	return &UserNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// UserInfoUpdatedNotification 发送用户信息更新通知
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

// UserStatusChangeNotification 发送用户状态变更通知
func (u *UserNotificationSender) UserStatusChangeNotification(ctx context.Context, tips *sdkws.UserStatusChangeTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserStatusChangeNotification, tips)
}

// UserCommandAddNotification 发送用户命令添加通知
func (u *UserNotificationSender) UserCommandAddNotification(ctx context.Context, tips *sdkws.UserCommandAddTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandAddNotification, tips)
}

// UserCommandDeleteNotification 发送用户命令删除通知
func (u *UserNotificationSender) UserCommandDeleteNotification(ctx context.Context, tips *sdkws.UserCommandDeleteTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandDeleteNotification, tips)
}

// UserCommandUpdateNotification 发送用户命令更新通知
func (u *UserNotificationSender) UserCommandUpdateNotification(ctx context.Context, tips *sdkws.UserCommandUpdateTips) {
	u.Notification(ctx, tips.FromUserID, tips.ToUserID, constant.UserCommandUpdateNotification, tips)
}
