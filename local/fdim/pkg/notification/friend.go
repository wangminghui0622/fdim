package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/database"
	"fdim/pkg/log"
	"fdim/pkg/mcontext"
	"fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// FriendNotificationSender 好友通知发送器
type FriendNotificationSender struct {
	*NotificationSender
	friendDB *database.FriendDatabase
}

// NewFriendNotificationSender 创建好友通知发送器
func NewFriendNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *FriendNotificationSender {
	return &FriendNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// SetFriendDB 设置好友数据库（用于获取好友申请信息）
func (f *FriendNotificationSender) SetFriendDB(db *database.FriendDatabase) {
	f.friendDB = db
}

// FriendApplicationAddNotification 发送好友申请通知
func (f *FriendNotificationSender) FriendApplicationAddNotification(ctx context.Context, fromUserID, toUserID string) {
	log.ZInfo(ctx, "[FriendNotification] FriendApplicationAddNotification called", "fromUserID", fromUserID, "toUserID", toUserID, "contentType", constant.FriendApplicationNotification)
	tips := &sdkws.FriendApplicationTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
		},
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendApplicationNotification, tips)
}

// FriendApplicationAgreedNotification 发送好友申请同意通知
// 参考官方实现：会获取好友申请信息（包含reqMsg）并发送，这样双方都能在消息列表看到对方
func (f *FriendNotificationSender) FriendApplicationAgreedNotification(ctx context.Context, fromUserID, toUserID, handleMsg string) {
	var request *sdkws.FriendRequest
	
	// 尝试获取好友申请信息（包含reqMsg验证消息）
	if f.friendDB != nil {
		requests, err := f.friendDB.FindBothFriendRequests(ctx, fromUserID, toUserID)
		if err != nil {
			log.ZError(ctx, "FriendApplicationAgreedNotification get friend request failed", err, "fromUserID", fromUserID, "toUserID", toUserID)
		} else {
			// 找到从fromUserID到toUserID的申请
			for _, req := range requests {
				if req.FromUserID == fromUserID && req.ToUserID == toUserID {
					request = &sdkws.FriendRequest{
						FromUserID:    req.FromUserID,
						ToUserID:      req.ToUserID,
						HandleResult:  req.HandleResult,
						ReqMsg:        req.ReqMsg, // 这是关键：申请时的验证消息
						CreateTime:    req.CreateTime.UnixMilli(),
						HandlerUserID: req.HandlerUserID,
						HandleMsg:     req.HandleMsg,
						HandleTime:    req.HandleTime.UnixMilli(),
						Ex:            req.Ex,
					}
					break
				}
			}
		}
	}
	
	tips := &sdkws.FriendApplicationApprovedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
		},
		HandleMsg: handleMsg,
		Request:   request, // 包含完整的申请信息，包括reqMsg
	}
	f.Notification(ctx, toUserID, fromUserID, constant.FriendApplicationApprovedNotification, tips)
}

// FriendApplicationRejectedNotification 发送好友申请拒绝通知
func (f *FriendNotificationSender) FriendApplicationRejectedNotification(ctx context.Context, fromUserID, toUserID, handleMsg string) {
	tips := &sdkws.FriendApplicationRejectedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
		},
		HandleMsg: handleMsg,
	}
	f.Notification(ctx, toUserID, fromUserID, constant.FriendApplicationRejectedNotification, tips)
}

// FriendAddedNotification 发送好友添加成功通知（双方都收到）
func (f *FriendNotificationSender) FriendAddedNotification(ctx context.Context, fromUserID, toUserID string) {
	tips := &sdkws.FriendAddedTips{
		Friend: &sdkws.FriendInfo{},
		OpUser: &sdkws.PublicUserInfo{UserID: fromUserID},
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendAddedNotification, tips)
}

// FriendDeletedNotification 发送好友删除通知
func (f *FriendNotificationSender) FriendDeletedNotification(ctx context.Context, ownerUserID, friendUserID string) {
	tips := &sdkws.FriendDeletedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   friendUserID,
		},
	}
	f.Notification(ctx, ownerUserID, friendUserID, constant.FriendDeletedNotification, tips)
}

// FriendRemarkSetNotification 发送好友备注设置通知
func (f *FriendNotificationSender) FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	tips := &sdkws.FriendInfoChangedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
		},
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendRemarkSetNotification, tips)
}

// FriendsInfoUpdateNotification 发送好友信息更新通知
func (f *FriendNotificationSender) FriendsInfoUpdateNotification(ctx context.Context, ownerUserID string, friendIDs []string) {
	tips := &sdkws.FriendsInfoUpdateTips{
		FromToUserID: &sdkws.FromToUserID{
			ToUserID: ownerUserID,
		},
		FriendIDs: friendIDs,
	}
	f.Notification(ctx, ownerUserID, ownerUserID, constant.FriendsInfoUpdateNotification, tips)
}

// BlackAddedNotification 发送黑名单添加通知
func (f *FriendNotificationSender) BlackAddedNotification(ctx context.Context, ownerUserID, blackUserID string) {
	tips := &sdkws.BlackAddedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   blackUserID,
		},
	}
	f.Notification(ctx, ownerUserID, blackUserID, constant.BlackAddedNotification, tips)
}

// BlackDeletedNotification 发送黑名单删除通知
func (f *FriendNotificationSender) BlackDeletedNotification(ctx context.Context, ownerUserID, blackUserID string) {
	tips := &sdkws.BlackDeletedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   blackUserID,
		},
	}
	f.Notification(ctx, ownerUserID, blackUserID, constant.BlackDeletedNotification, tips)
}

// FriendInfoUpdatedNotification 发送好友信息更新通知
func (f *FriendNotificationSender) FriendInfoUpdatedNotification(ctx context.Context, changedUserID, needNotifiedUserID string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = changedUserID
	}
	tips := &sdkws.UserInfoUpdatedTips{
		UserID: changedUserID,
	}
	f.Notification(ctx, opUserID, needNotifiedUserID, constant.FriendInfoUpdatedNotification, tips)
}
