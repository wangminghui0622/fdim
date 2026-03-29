package notification

import (
	"context"

	"fdim/pkg/config"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	protoconstant "fdim/protocol/constant"
	"fdim/protocol/sdkws"
)

// GroupNotificationSender Ⱥ֪ͨ
type GroupNotificationSender struct {
	*NotificationSender
}

// NewGroupNotificationSender Ⱥ֪ͨ
func NewGroupNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *GroupNotificationSender {
	return &GroupNotificationSender{
		NotificationSender: NewNotificationSender(conf, opts...),
	}
}

// GroupCreatedNotification 群创建通知
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
	// 发送给群（所有成员都会收到）
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupCreatedNotification, constant.SuperGroupChatType, tips)
}

// GroupInfoSetNotification ȺϢ֪ͨ
func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupInfoSetTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// ͸Ⱥ飨Ⱥͣ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupInfoSetNotification, constant.SuperGroupChatType, tips)
}

// GroupOwnerTransferredNotification Ⱥת֪ͨ
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
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, oldOwnerUserID, groupID, protoconstant.GroupOwnerTransferredNotification, constant.SuperGroupChatType, tips)
}

// MemberKickedNotification ͳԱ֪߳ͨ
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
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberKickedNotification, constant.SuperGroupChatType, tips)
}

// MemberQuitNotification ͳԱ˳֪ͨ
func (g *GroupNotificationSender) MemberQuitNotification(ctx context.Context, groupID string, quitUserID string) {
	tips := &sdkws.MemberQuitTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		QuitUser: &sdkws.GroupMemberFullInfo{
			UserID: quitUserID,
		},
	}
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, quitUserID, groupID, protoconstant.MemberQuitNotification, constant.SuperGroupChatType, tips)
}

// MemberEnterNotification ͳԱ֪ͨ
func (g *GroupNotificationSender) MemberEnterNotification(ctx context.Context, groupID string, inviterUserID string, newMemberUserIDs []string) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		opUserID = inviterUserID
	}
	// ע⣺MemberEnterTips.EntrantUser ǵ󣬲Ƭ
	// жԱҪΪÿԱ͵֪ͨ
	if len(newMemberUserIDs) == 0 {
		return
	}
	// Ϊÿ³Ա͵֪ͨ
	for _, userID := range newMemberUserIDs {
		tips := &sdkws.MemberEnterTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			EntrantUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberEnterNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupDismissedNotification Ⱥɢ֪ͨ
func (g *GroupNotificationSender) GroupDismissedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupDismissedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupDismissedNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberInfoSetNotification ȺԱϢ֪ͨ
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
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberInfoSetNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberSetToAdminNotification ȺԱΪԱ֪ͨ
// ע⣺proto û GroupMemberSetToAdminTipsʹ GroupMemberInfoSetTips ΪÿԱ֪ͨ
func (g *GroupNotificationSender) GroupMemberSetToAdminNotification(ctx context.Context, groupID string, opUserID string, newAdminUserIDs []string) {
	// Ϊÿ¹Ա͵֪ͨ
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
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberSetToAdminNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMemberSetToOrdinaryUserNotification ȺԱΪͨû֪ͨ
// ע⣺proto û GroupMemberSetToOrdinaryUserTipsʹ GroupMemberInfoSetTips ΪÿԱ֪ͨ
func (g *GroupNotificationSender) GroupMemberSetToOrdinaryUserNotification(ctx context.Context, groupID string, opUserID string, setToOrdinaryUserIDs []string) {
	// ΪÿΪͨûĳԱ͵֪ͨ
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
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberSetToOrdinaryUserNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupApplicationAcceptedNotification Ⱥ֪ͨ
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
	// ͸??
	g.Notification(ctx, opUserID, applicantUserID, protoconstant.GroupApplicationAcceptedNotification, tips)
}

// GroupApplicationRejectedNotification Ⱥܾ֪ͨ
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
	// ͸??
	g.Notification(ctx, opUserID, applicantUserID, protoconstant.GroupApplicationRejectedNotification, tips)
}

// JoinGroupApplicationNotification ͼȺ֪ͨ
func (g *GroupNotificationSender) JoinGroupApplicationNotification(ctx context.Context, groupID string, applicantUserID string) {
	tips := &sdkws.JoinGroupApplicationTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		Applicant: &sdkws.PublicUserInfo{
			UserID: applicantUserID,
		},
	}
	// ͸Ⱥ飨֪ͨԱ
	g.NotificationWithSessionType(ctx, applicantUserID, groupID, protoconstant.JoinGroupApplicationNotification, constant.SuperGroupChatType, tips)
}

// GroupApplicationAgreeMemberEnterNotification ȺͬԱ֪ͨ
// ע⣺ʹ GroupApplicationAcceptedTipsΪ GroupApplicationAgreedTips 
func (g *GroupNotificationSender) GroupApplicationAgreeMemberEnterNotification(ctx context.Context, groupID string, opUserID string, newMemberUserIDs []string) {
	// ʹ MemberEnterNotification Ϊÿ³Ա֪ͨ
	for _, userID := range newMemberUserIDs {
		tips := &sdkws.MemberEnterTips{
			Group: &sdkws.GroupInfo{
				GroupID: groupID,
			},
			EntrantUser: &sdkws.GroupMemberFullInfo{
				UserID: userID,
			},
		}
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.MemberEnterNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMutedNotification Ⱥ֪ͨ
func (g *GroupNotificationSender) GroupMutedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupMutedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMutedNotification, constant.SuperGroupChatType, tips)
}

// GroupCancelMutedNotification Ⱥȡ֪ͨ
func (g *GroupNotificationSender) GroupCancelMutedNotification(ctx context.Context, groupID string, opUserID string) {
	tips := &sdkws.GroupCancelMutedTips{
		Group: &sdkws.GroupInfo{
			GroupID: groupID,
		},
		OpUser: &sdkws.GroupMemberFullInfo{
			UserID: opUserID,
		},
	}
	// ͸Ⱥ
	g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupCancelMutedNotification, constant.SuperGroupChatType, tips)
}

// GroupMemberMutedNotification ȺԱ֪ͨ
// ע⣺GroupMemberMutedTips ֻһ mutedUser ֶΣҪΪÿԱ͵֪ͨ
func (g *GroupNotificationSender) GroupMemberMutedNotification(ctx context.Context, groupID string, opUserID string, mutedUserIDs []string) {
	// ΪÿԵĳԱ͵֪ͨ
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
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberMutedNotification, constant.SuperGroupChatType, tips)
	}
}

// GroupMemberCancelMutedNotification ȺԱȡ֪ͨ
// ע⣺GroupMemberCancelMutedTips ֻһ mutedUser ֶΣҪΪÿԱ͵֪ͨ
func (g *GroupNotificationSender) GroupMemberCancelMutedNotification(ctx context.Context, groupID string, opUserID string, unmutedUserIDs []string) {
	// ΪÿȡԵĳԱ͵֪ͨ
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
		// ͸Ⱥ
		g.NotificationWithSessionType(ctx, opUserID, groupID, protoconstant.GroupMemberCancelMutedNotification, constant.SuperGroupChatType, tips)
	}
}
