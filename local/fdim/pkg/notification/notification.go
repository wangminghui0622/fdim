package notification

import (
	"fdim/pkg/constant"
	protoconstant "fdim/protocol/constant"
)

// getSessionType 根据 contentType 获取 sessionType
func getSessionType(contentType int32) int32 {
	// 好友相关通知使用单聊
	if contentType == protoconstant.FriendApplicationNotification ||
		contentType == protoconstant.FriendApplicationApprovedNotification ||
		contentType == protoconstant.FriendApplicationRejectedNotification ||
		contentType == protoconstant.FriendAddedNotification ||
		contentType == protoconstant.FriendDeletedNotification ||
		contentType == protoconstant.FriendRemarkSetNotification ||
		contentType == protoconstant.BlackAddedNotification ||
		contentType == protoconstant.BlackDeletedNotification ||
		contentType == protoconstant.FriendInfoUpdatedNotification ||
		contentType == protoconstant.FriendsInfoUpdateNotification {
		return constant.SingleChatType
	}
	// 群组相关通知使用群聊
	if contentType == protoconstant.GroupCreatedNotification ||
		contentType == protoconstant.GroupInfoSetNotification ||
		contentType == protoconstant.GroupOwnerTransferredNotification ||
		contentType == protoconstant.MemberKickedNotification ||
		contentType == protoconstant.MemberQuitNotification ||
		contentType == protoconstant.MemberInvitedNotification ||
		contentType == protoconstant.MemberEnterNotification ||
		contentType == protoconstant.GroupDismissedNotification ||
		contentType == protoconstant.GroupMemberMutedNotification ||
		contentType == protoconstant.GroupMemberCancelMutedNotification ||
		contentType == protoconstant.GroupMutedNotification ||
		contentType == protoconstant.GroupCancelMutedNotification ||
		contentType == protoconstant.GroupMemberInfoSetNotification ||
		contentType == protoconstant.GroupMemberSetToAdminNotification ||
		contentType == protoconstant.GroupMemberSetToOrdinaryUserNotification ||
		contentType == protoconstant.GroupApplicationAcceptedNotification ||
		contentType == protoconstant.GroupApplicationRejectedNotification ||
		contentType == protoconstant.JoinGroupApplicationNotification {
		return constant.SuperGroupChatType
	}
	// 用户相关通知使用单聊
	if contentType == protoconstant.UserInfoUpdatedNotification ||
		contentType == protoconstant.UserStatusChangeNotification ||
		contentType == protoconstant.UserCommandAddNotification ||
		contentType == protoconstant.UserCommandDeleteNotification ||
		contentType == protoconstant.UserCommandUpdateNotification {
		return constant.SingleChatType
	}
	// 会话相关通知使用单聊（与官方一致）
	if contentType == protoconstant.ConversationChangeNotification ||
		contentType == protoconstant.ConversationPrivateChatNotification ||
		contentType == protoconstant.ConversationUnreadNotification ||
		contentType == protoconstant.ClearConversationNotification ||
		contentType == protoconstant.ConversationDeleteNotification {
		return constant.SingleChatType
	}
	// 默认使用单聊
	return constant.SingleChatType
}
