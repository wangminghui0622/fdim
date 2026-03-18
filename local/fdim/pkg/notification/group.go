package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	protoconstant "fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// GroupNotificationSender 群组通知发送器
type GroupNotificationSender struct {
	*NotificationSender
}

// NewGroupNotificationSender 创建群组通知发送器
func NewGroupNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *GroupNotificationSender {
	return &GroupNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// GroupCreatedNotification 发送群组创建通知
func (g *GroupNotificationSender) GroupCreatedNotification(ctx context.Context, groupID string, ownerUserID string, memberUserIDs []string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = ownerUserID
	}
	tips := &sdkws.GroupCreatedTips{
		Group: &sdkws.GroupInfo{
			GroupID:     groupID,
			OwnerUserID: ownerUserID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// 发送给所有群成员
	for _, memberID := range memberUserIDs {
		g.Notification(ctx, opUserID, memberID, protoconstant.GroupCreatedNotification, tips)
	}
}

// GroupInfoSetNotification 发送群组信息设置通知
func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupInfoSetTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// 发送给群组（群聊类型）
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupInfoSetNotification, constant.SuperGroupChatType, tips)
}

// GroupOwnerTransferredNotification 发送群主转让通知
func (g *GroupNotificationSender) GroupOwnerTransferredNotification(ctx context.Context, groupID string, oldOwnerUserID string, newOwnerUserID string) {
	tips := &sdkws.GroupOwnerTransferredTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: oldOwnerUserID,
		},
		NewGroupOwner: &sdkws.GroupMemberFullInfo{
			UserID: newOwnerUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, oldOwnerUserID, groupID, protoconstant.GroupOwnerTransferredNotification, constant.SuperGroupChatType, tips)
}

// MemberKickedNotification 发送成员被踢出通知
func (g *GroupNotificationSender) MemberKickedNotification(ctx context.Context, groupID string, opUserID string, kickedUserIDs []string) {
	tips := &sdkws.MemberKickedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
		KickedUserList: make([]*sdkws.GroupMemberFullInfo, 0, len(kickedUserIDs)),
	}
	for _, userID := range kickedUserIDs {
		tips.KickedUserList = append(tips.KickedUserList, &sdkws.GroupMemberFullInfo{
			UserID: userID,
		})
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberKickedNotification, constant.SuperGroupChatType, tips)
}

// MemberQuitNotification 发送成员退出通知
func (g *GroupNotificationSender) MemberQuitNotification(ctx context.Context, groupID string, quitUserID string) {
	tips := &sdkws.MemberQuitTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		QuitUser: &sdkws.GroupMemberFullInfo{
			UserID: quitUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, quitUserID, groupID, protoconstant.MemberQuitNotification, constant.SuperGroupChatType, tips)
}

// MemberEnterNotification 发送成员加入通知
func (g *GroupNotificationSender) MemberEnterNotification(ctx context.Context, groupID string, inviterUserID string, newMemberUserIDs []string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = inviterUserID
	}
	// 注意：MemberEnterTips.EntrantUser 是单个对象，不是切片
	// 如果有多个成员，需要为每个成员发送单独的通知
	if len(newMemberUserIDs) == 0 {
		return
	}
	// 为每个新成员发送单独的通知
	for _, userID := range newMemberUserIDs {
		tips := &sdkws.MemberEnterTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			EntrantUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberEnterNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupDismissedNotification 发送群组解散通知
func (g *GroupNotificationSender) GroupDismissedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupDismissedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupDismissedNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberInfoSetNotification 发送群成员信息设置通知
func (g *GroupNotificationSender) GroupMemberInfoSetNotification(ctx context.Context, groupID string, opUserID string, changedUserID string) {
	tips := &sdkws.GroupMemberInfoSetTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
		ChangedUser: &sdkws.GroupMemberFullInfo{
			UserID: changedUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberInfoSetNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberSetToAdminNotification 发送群成员设置为管理员通知
// 注意：proto 中没有 GroupMemberSetToAdminTips，使用 GroupMemberInfoSetTips 为每个成员发送通知
func (g *GroupNotificationSender) GroupMemberSetToAdminNotification(ctx context.Context, groupID string, opUserID string, newAdminUserIDs []string) {
	// 为每个新管理员发送单独的通知
	for _, userID := range newAdminUserIDs {
		tips := &sdkws.GroupMemberInfoSetTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			OpUser: &sdkws.GroupMemberFullInfo{
				UserID: opUserID,
			},
			ChangedUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberSetToAdminNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMemberSetToOrdinaryUserNotification 发送群成员设置为普通用户通知
// 注意：proto 中没有 GroupMemberSetToOrdinaryUserTips，使用 GroupMemberInfoSetTips 为每个成员发送通知
func (g *GroupNotificationSender) GroupMemberSetToOrdinaryUserNotification(ctx context.Context, groupID string, opUserID string, setToOrdinaryUserIDs []string) {
	// 为每个被设置为普通用户的成员发送单独的通知
	for _, userID := range setToOrdinaryUserIDs {
		tips := &sdkws.GroupMemberInfoSetTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			OpUser: &sdkws.GroupMemberFullInfo{
				UserID: opUserID,
			},
			ChangedUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberSetToOrdinaryUserNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupApplicationAcceptedNotification 发送群组申请接受通知
func (g *GroupNotificationSender) GroupApplicationAcceptedNotification(ctx context.Context, groupID string, opUserID string, applicantUserID string) {
	tips := &sdkws.GroupApplicationAcceptedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
		HandleMsg: "",
	}
	// 发送给申请�?
	g.Notification(ctx, opUserID, applicantUserID, protoconstant.GroupApplicationAcceptedNotification, tips)
}

// GroupApplicationRejectedNotification 发送群组申请拒绝通知
func (g *GroupNotificationSender) GroupApplicationRejectedNotification(ctx context.Context, groupID string, opUserID string, applicantUserID string, handleMsg string) {
	tips := &sdkws.GroupApplicationRejectedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
		HandleMsg: handleMsg,
	}
	// 发送给申请�?
	g.Notification(ctx, opUserID, applicantUserID, protoconstant.GroupApplicationRejectedNotification, tips)
}

// JoinGroupApplicationNotification 发送加入群组申请通知
func (g *GroupNotificationSender) JoinGroupApplicationNotification(ctx context.Context, groupID string, applicantUserID string) {
	tips := &sdkws.JoinGroupApplicationTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		Applicant: &sdkws.PublicUserInfo{
			UserID: applicantUserID,
		},
	}
	// 发送给群组（通知管理员）
	g.NotificationWithSessionType(ctx, applicantUserID, groupID, protoconstant.JoinGroupApplicationNotification, constant.SuperGroupChatType, tips)
}

// GroupApplicationAgreeMemberEnterNotification 发送群组申请同意后成员加入通知
// 注意：使用 GroupApplicationAcceptedTips，因为 GroupApplicationAgreedTips 不存在
func (g *GroupNotificationSender) GroupApplicationAgreeMemberEnterNotification(ctx context.Context, groupID string, opUserID string, newMemberUserIDs []string) {
	// 使用 MemberEnterNotification 为每个新成员发送通知
	for _, userID := range newMemberUserIDs {
		tips := &sdkws.MemberEnterTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			EntrantUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberEnterNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMutedNotification 发送群组禁言通知
func (g *GroupNotificationSender) GroupMutedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupMutedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMutedNotification, constant.SuperGroupChatType, tips)
}

// GroupCancelMutedNotification 发送群组取消禁言通知
func (g *GroupNotificationSender) GroupCancelMutedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupCancelMutedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// 发送给群组
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupCancelMutedNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberMutedNotification 发送群成员禁言通知
// 注意：GroupMemberMutedTips 只有一个 mutedUser 字段，需要为每个成员发送单独的通知
func (g *GroupNotificationSender) GroupMemberMutedNotification(ctx context.Context, groupID string, opUserID string, mutedUserIDs []string) {
	// 为每个被禁言的成员发送单独的通知
	for _, userID := range mutedUserIDs {
		tips := &sdkws.GroupMemberMutedTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			OpUser: &sdkws.GroupMemberFullInfo{
				UserID: opUserID,
			},
			MutedUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberMutedNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMemberCancelMutedNotification 发送群成员取消禁言通知
// 注意：GroupMemberCancelMutedTips 只有一个 mutedUser 字段，需要为每个成员发送单独的通知
func (g *GroupNotificationSender) GroupMemberCancelMutedNotification(ctx context.Context, groupID string, opUserID string, unmutedUserIDs []string) {
	// 为每个被取消禁言的成员发送单独的通知
	for _, userID := range unmutedUserIDs {
		tips := &sdkws.GroupMemberCancelMutedTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			OpUser: &sdkws.GroupMemberFullInfo{
				UserID: opUserID,
			},
			MutedUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// 发送给群组
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberCancelMutedNotification, constant.SuperGroupChatType, tips)
	}
}
