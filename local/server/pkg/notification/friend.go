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

// FriendNotificationSender ֪ͨ
type FriendNotificationSender struct {
	*NotificationSender
	friendDB *database.FriendDatabase
}

// NewFriendNotificationSender ֪ͨ
func NewFriendNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *FriendNotificationSender {
	return &FriendNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// SetFriendDB úݿ⣨ڻȡϢ
func (f *FriendNotificationSender) SetFriendDB(db *database.FriendDatabase) {
	f.friendDB = db
}

// FriendApplicationAddNotification ͺ֪ͨ
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

// FriendApplicationAgreedNotification ͺ֪ͬͨ
// οٷʵ֣ȡϢreqMsgͣ˫ϢбԷ
func (f *FriendNotificationSender) FriendApplicationAgreedNotification(ctx context.Context, fromUserID, toUserID, handleMsg string) {
	var request *sdkws.FriendRequest
	
	// ԻȡϢreqMsg֤Ϣ
	if f.friendDB != nil {
		requests, err := f.friendDB.FindBothFriendRequests(ctx, fromUserID, toUserID)
		if err != nil {
			log.ZError(ctx, "FriendApplicationAgreedNotification get friend request failed", err, "fromUserID", fromUserID, "toUserID", toUserID)
		} else {
			// ҵfromUserIDtoUserID
			for _, req := range requests {
				if req.FromUserID == fromUserID && req.ToUserID == toUserID {
					request = &sdkws.FriendRequest{
						FromUserID:    req.FromUserID,
						ToUserID:      req.ToUserID,
						HandleResult:  req.HandleResult,
						ReqMsg:        req.ReqMsg, // ǹؼʱ֤Ϣ
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
		Request:   request, // ϢreqMsg
	}
	f.Notification(ctx, toUserID, fromUserID, constant.FriendApplicationApprovedNotification, tips)
}

// FriendApplicationRejectedNotification ͺܾ֪ͨ
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

// FriendAddedNotification ͺӳɹ֪ͨ˫յ
func (f *FriendNotificationSender) FriendAddedNotification(ctx context.Context, fromUserID, toUserID string) {
	tips := &sdkws.FriendAddedTips{
		Friend: &sdkws.FriendInfo{},
		OpUser: &sdkws.PublicUserInfo{UserID: fromUserID},
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendAddedNotification, tips)
}

// FriendDeletedNotification ͺɾ֪ͨ
func (f *FriendNotificationSender) FriendDeletedNotification(ctx context.Context, ownerUserID, friendUserID string) {
	tips := &sdkws.FriendDeletedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   friendUserID,
		},
	}
	f.Notification(ctx, ownerUserID, friendUserID, constant.FriendDeletedNotification, tips)
}

// FriendRemarkSetNotification ͺѱע֪ͨ
func (f *FriendNotificationSender) FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	tips := &sdkws.FriendInfoChangedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
		},
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendRemarkSetNotification, tips)
}

// FriendsInfoUpdateNotification ͺϢ֪ͨ
func (f *FriendNotificationSender) FriendsInfoUpdateNotification(ctx context.Context, ownerUserID string, friendIDs []string) {
	tips := &sdkws.FriendsInfoUpdateTips{
		FromToUserID: &sdkws.FromToUserID{
			ToUserID: ownerUserID,
		},
		FriendIDs: friendIDs,
	}
	f.Notification(ctx, ownerUserID, ownerUserID, constant.FriendsInfoUpdateNotification, tips)
}

// BlackAddedNotification ͺ֪ͨ
func (f *FriendNotificationSender) BlackAddedNotification(ctx context.Context, ownerUserID, blackUserID string) {
	tips := &sdkws.BlackAddedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   blackUserID,
		},
	}
	f.Notification(ctx, ownerUserID, blackUserID, constant.BlackAddedNotification, tips)
}

// BlackDeletedNotification ͺɾ֪ͨ
func (f *FriendNotificationSender) BlackDeletedNotification(ctx context.Context, ownerUserID, blackUserID string) {
	tips := &sdkws.BlackDeletedTips{
		FromToUserID: &sdkws.FromToUserID{
			FromUserID: ownerUserID,
			ToUserID:   blackUserID,
		},
	}
	f.Notification(ctx, ownerUserID, blackUserID, constant.BlackDeletedNotification, tips)
}

// FriendInfoUpdatedNotification ͺϢ֪ͨ
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
