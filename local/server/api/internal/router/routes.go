package router

import (
	"fdim/api/internal/handler"
	"fdim/api/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

// RegisterHandlers 将所有 HTTP 路由注册到 server
func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	// User 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/user_register",
		Handler: handler.UserRegisterHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/update_user_info",
		Handler: handler.UpdateUserInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/update_user_info_ex",
		Handler: handler.UpdateUserInfoExHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_users_info",
		Handler: handler.GetUsersInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/account_check",
		Handler: handler.AccountCheckHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/set_global_msg_recv_opt",
		Handler: handler.SetGlobalRecvMessageOptHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_all_users_uid",
		Handler: handler.GetAllUsersIDHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_users",
		Handler: handler.GetUsersHandler(ctx),
	})

	// user_register_count 移到 /statistics 组（与官方一致）

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_users_status",
		Handler: handler.GetUserStatusHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_subscribe_users_status",
		Handler: handler.GetSubscribeUsersStatusHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/subscribe_users_status",
		Handler: handler.SubscribeUsersStatusHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/process_user_command_add",
		Handler: handler.ProcessUserCommandAddHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/process_user_command_delete",
		Handler: handler.ProcessUserCommandDeleteHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/process_user_command_update",
		Handler: handler.ProcessUserCommandUpdateHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/process_user_command_get",
		Handler: handler.ProcessUserCommandGetHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/process_user_command_get_all",
		Handler: handler.ProcessUserCommandGetAllHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/add_notification_account",
		Handler: handler.AddNotificationAccountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/update_notification_account",
		Handler: handler.UpdateNotificationAccountInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/search_notification_account",
		Handler: handler.SearchNotificationAccountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_user_client_config",
		Handler: handler.GetUserClientConfigHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/set_user_client_config",
		Handler: handler.SetUserClientConfigHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/del_user_client_config",
		Handler: handler.DelUserClientConfigHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/page_user_client_config",
		Handler: handler.PageUserClientConfigHandler(ctx),
	})
	//好友在线状态查询
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_users_online_status",
		Handler: handler.GetUsersOnlineStatusHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/get_users_online_token_detail",
		Handler: handler.GetUsersOnlineTokenDetailHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/search/full",
		Handler: handler.SearchUserFullInfoHandler(ctx),
	})

	// Friend 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/add_friend",
		Handler: handler.ApplyToAddFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/add_friend_response",
		Handler: handler.RespondFriendApplyHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/delete_friend",
		Handler: handler.DeleteFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_friend_apply_list",
		Handler: handler.GetFriendApplyListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_designated_friend_apply",
		Handler: handler.GetDesignatedFriendsApplyHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_self_friend_apply_list",
		Handler: handler.GetSelfApplyListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_friend_list",
		Handler: handler.GetFriendListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_designated_friends",
		Handler: handler.GetDesignatedFriendsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/set_friend_remark",
		Handler: handler.SetFriendRemarkHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/add_black",
		Handler: handler.AddBlackHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_black_list",
		Handler: handler.GetPaginationBlacksHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_specified_blacks",
		Handler: handler.GetSpecifiedBlacksHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/remove_black",
		Handler: handler.RemoveBlackHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_incremental_blacks",
		Handler: handler.GetIncrementalBlacksHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/import_friend",
		Handler: handler.ImportFriendsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/is_friend",
		Handler: handler.IsFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_friend_id",
		Handler: handler.GetFriendIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_specified_friends_info",
		Handler: handler.GetSpecifiedFriendsInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/update_friends",
		Handler: handler.UpdateFriendsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_incremental_friends",
		Handler: handler.GetIncrementalFriendsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_full_friend_user_ids",
		Handler: handler.GetFullFriendUserIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/get_self_unhandled_apply_count",
		Handler: handler.GetSelfUnhandledApplyCountHandler(ctx),
	})

	// Group 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/create_group",
		Handler: handler.CreateGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/set_group_info",
		Handler: handler.SetGroupInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/set_group_info_ex",
		Handler: handler.SetGroupInfoExHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/join_group",
		Handler: handler.JoinGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/quit_group",
		Handler: handler.QuitGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/group_application_response",
		Handler: handler.ApplicationGroupResponseHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/transfer_group",
		Handler: handler.TransferGroupOwnerHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_recv_group_applicationList",
		Handler: handler.GetRecvGroupApplicationListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_user_req_group_applicationList",
		Handler: handler.GetUserReqGroupApplicationListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_users_req_application_list",
		Handler: handler.GetGroupUsersReqApplicationListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_specified_user_group_request_info",
		Handler: handler.GetSpecifiedUserGroupRequestInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_groups_info",
		Handler: handler.GetGroupsInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/kick_group",
		Handler: handler.KickGroupMemberHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_members_info",
		Handler: handler.GetGroupMembersInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_member_list",
		Handler: handler.GetGroupMemberListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/invite_user_to_group",
		Handler: handler.InviteUserToGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_joined_group_list",
		Handler: handler.GetJoinedGroupListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/dismiss_group",
		Handler: handler.DismissGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/mute_group_member",
		Handler: handler.MuteGroupMemberHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/cancel_mute_group_member",
		Handler: handler.CancelMuteGroupMemberHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/mute_group",
		Handler: handler.MuteGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/cancel_mute_group",
		Handler: handler.CancelMuteGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/set_group_member_info",
		Handler: handler.SetGroupMemberInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_abstract_info",
		Handler: handler.GetGroupAbstractInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_groups",
		Handler: handler.GetGroupsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_member_user_id",
		Handler: handler.GetGroupMemberUserIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_incremental_join_groups",
		Handler: handler.GetIncrementalJoinGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_incremental_group_members",
		Handler: handler.GetIncrementalGroupMemberHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_incremental_group_members_batch",
		Handler: handler.GetIncrementalGroupMemberBatchHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_full_group_member_user_ids",
		Handler: handler.GetFullGroupMemberUserIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_full_join_group_ids",
		Handler: handler.GetFullJoinGroupIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/group/get_group_application_unhandled_count",
		Handler: handler.GetGroupApplicationUnhandledCountHandler(ctx),
	})

	// Auth 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/auth/get_admin_token",
		Handler: handler.GetAdminTokenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/auth/get_user_token",
		Handler: handler.GetUserTokenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/auth/parse_token",
		Handler: handler.ParseTokenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/auth/force_logout",
		Handler: handler.ForceLogoutHandler(ctx),
	})

	// Conversation 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_sorted_conversation_list",
		Handler: handler.GetSortedConversationListHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_all_conversations",
		Handler: handler.GetAllConversationsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_conversation",
		Handler: handler.GetConversationHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_conversations",
		Handler: handler.GetConversationsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/set_conversations",
		Handler: handler.SetConversationsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_full_conversation_ids",
		Handler: handler.GetFullOwnerConversationIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_incremental_conversations",
		Handler: handler.GetIncrementalConversationHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_owner_conversation",
		Handler: handler.GetOwnerConversationHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_not_notify_conversation_ids",
		Handler: handler.GetNotNotifyConversationIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/get_pinned_conversation_ids",
		Handler: handler.GetPinnedConversationIDsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/delete_conversations",
		Handler: handler.DeleteConversationsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/conversation/update_conversations_by_user",
		Handler: handler.UpdateConversationsByUserHandler(ctx),
	})

	// Msg 相关路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/newest_seq",
		Handler: handler.GetSeqHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/search_msg",
		Handler: handler.SearchMsgHandler(ctx),
	})
	//发送消息
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/send_msg",
		Handler: handler.SendMessageHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/send_business_notification",
		Handler: handler.SendBusinessNotificationHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/pull_msg_by_seq",
		Handler: handler.PullMsgBySeqsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/revoke_msg",
		Handler: handler.RevokeMsgHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/mark_msgs_as_read",
		Handler: handler.MarkMsgsAsReadHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/mark_conversation_as_read",
		Handler: handler.MarkConversationAsReadHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/get_conversations_has_read_and_max_seq",
		Handler: handler.GetConversationsHasReadAndMaxSeqHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/set_conversation_has_read_seq",
		Handler: handler.SetConversationHasReadSeqHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/clear_conversation_msg",
		Handler: handler.ClearConversationsMsgHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/user_clear_all_msg",
		Handler: handler.UserClearAllMsgHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/delete_msgs",
		Handler: handler.DeleteMsgsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/delete_msg_phsical_by_seq",
		Handler: handler.DeleteMsgPhysicalBySeqHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/delete_msg_physical",
		Handler: handler.DeleteMsgPhysicalHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/batch_send_msg",
		Handler: handler.BatchSendMsgHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/send_simple_msg",
		Handler: handler.SendSimpleMessageHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/check_msg_is_send_success",
		Handler: handler.CheckMsgIsSendSuccessHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/msg/get_server_time",
		Handler: handler.GetServerTimeHandler(ctx),
	})

	// Third 相关路由
	server.AddRoute(rest.Route{
		Method:  "GET",
		Path:    "/third/prometheus",
		Handler: handler.GetPrometheusHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/third/fcm_update_token",
		Handler: handler.FcmUpdateTokenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/third/set_app_badge",
		Handler: handler.SetAppBadgeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/third/logs/upload",
		Handler: handler.UploadLogsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/third/logs/delete",
		Handler: handler.DeleteLogsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/third/logs/search",
		Handler: handler.SearchLogsHandler(ctx),
	})

	// Object 相关路由（官方在 /object 根路径下，不在 /third/object 下）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/part_limit",
		Handler: handler.PartLimitHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/part_size",
		Handler: handler.PartSizeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/initiate_multipart_upload",
		Handler: handler.InitiateMultipartUploadHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/auth_sign",
		Handler: handler.AuthSignHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/complete_multipart_upload",
		Handler: handler.CompleteMultipartUploadHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/access_url",
		Handler: handler.AccessURLHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/initiate_form_data",
		Handler: handler.InitiateFormDataHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/object/complete_form_data",
		Handler: handler.CompleteFormDataHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "GET",
		Path:    "/object/:name",
		Handler: handler.ObjectRedirectHandler(ctx),
	})

	// Statistics 统计相关路由（与官方一致，在 /statistics 组下）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistics/user/register",
		Handler: handler.UserRegisterCountHandler(ctx),
	})

	// Favorite 收藏相关路由（自定义扩展）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/favorite/add",
		Handler: handler.AddFavoriteHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/favorite/delete",
		Handler: handler.DeleteFavoriteHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/favorite/list",
		Handler: handler.GetFavoriteListHandler(ctx),
	})

	// Account 登录/注册路由（通过 chat RPC 验证密码）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/login",
		Handler: handler.AccountLoginHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/register",
		Handler: handler.AccountRegisterHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/password/change",
		Handler: handler.ChangePasswordHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/password/reset",
		Handler: handler.ResetPasswordHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/friend/search",
		Handler: handler.SearchFriendHandler(ctx),
	})

	// 验证码路由（对齐官方 chat server）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/code/send",
		Handler: handler.SendVerifyCodeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/account/code/verify",
		Handler: handler.VerifyCodeHandler(ctx),
	})

	// Chat 层用户接口（对齐官方 chat server）
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/find/full",
		Handler: handler.FindUserFullInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/find/public",
		Handler: handler.FindUserPublicInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/search/public",
		Handler: handler.SearchUserPublicInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/rtc/get_token",
		Handler: handler.GetTokenForRTCHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/update",
		Handler: handler.ChatUpdateUserInfoHandler(ctx),
	})

	// 客户端配置
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/client_config/get",
		Handler: handler.GetClientConfigHandler(ctx),
	})

	// Applet 小程序
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/applet/find",
		Handler: handler.FindAppletHandler(ctx),
	})

	// Application 版本管理
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/application/latest_version",
		Handler: handler.LatestApplicationVersionHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/application/page_versions",
		Handler: handler.PageApplicationVersionHandler(ctx),
	})

	// OpenIM 回调
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/callback/open_im",
		Handler: handler.OpenIMCallbackHandler(ctx),
	})

	// Statistics 统计扩展路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistics/user/active",
		Handler: handler.GetActiveUserHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistics/group/create",
		Handler: handler.GroupCreateCountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistics/group/active",
		Handler: handler.GetActiveGroupHandler(ctx),
	})

	// JSSDK 路由
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/jssdk/get_conversations",
		Handler: handler.JSSdkGetConversationsHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/jssdk/get_active_conversations",
		Handler: handler.JSSdkGetActiveConversationsHandler(ctx),
	})

	// ========== Admin 管理路由 ==========

	// Admin 账号管理
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/login",
		Handler: handler.AdminLoginHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/update",
		Handler: handler.AdminUpdateInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/info",
		Handler: handler.AdminInfoHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/change_password",
		Handler: handler.ChangeAdminPasswordHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/add_admin",
		Handler: handler.AddAdminAccountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/add_user",
		Handler: handler.AddUserAccountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/del_admin",
		Handler: handler.DelAdminAccountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/search",
		Handler: handler.SearchAdminAccountHandler(ctx),
	})

	// Admin: 用户导入
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/import/json",
		Handler: handler.ImportUserByJsonHandler(ctx),
	})

	// Admin: 注册开关
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/allow_register/get",
		Handler: handler.GetAllowRegisterHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/user/allow_register/set",
		Handler: handler.SetAllowRegisterHandler(ctx),
	})

	// Admin: 默认好友
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/user/add",
		Handler: handler.AddDefaultFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/user/del",
		Handler: handler.DelDefaultFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/user/find",
		Handler: handler.FindDefaultFriendHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/user/search",
		Handler: handler.SearchDefaultFriendHandler(ctx),
	})

	// Admin: 默认群组
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/group/add",
		Handler: handler.AddDefaultGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/group/del",
		Handler: handler.DelDefaultGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/group/find",
		Handler: handler.FindDefaultGroupHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/default/group/search",
		Handler: handler.SearchDefaultGroupHandler(ctx),
	})

	// Admin: 邀请码
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/invitation_code/add",
		Handler: handler.AddInvitationCodeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/invitation_code/gen",
		Handler: handler.GenInvitationCodeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/invitation_code/del",
		Handler: handler.DelInvitationCodeHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/invitation_code/search",
		Handler: handler.SearchInvitationCodeHandler(ctx),
	})

	// Admin: IP 封禁
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/ip/add",
		Handler: handler.AddIPForbiddenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/ip/del",
		Handler: handler.DelIPForbiddenHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/ip/search",
		Handler: handler.SearchIPForbiddenHandler(ctx),
	})

	// Admin: 用户 IP 限制
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/user/add",
		Handler: handler.AddUserIPLimitLoginHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/user/del",
		Handler: handler.DelUserIPLimitLoginHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/forbidden/user/search",
		Handler: handler.SearchUserIPLimitLoginHandler(ctx),
	})

	// Admin: 小程序管理
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/applet/add",
		Handler: handler.AddAppletHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/applet/del",
		Handler: handler.DelAppletHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/applet/update",
		Handler: handler.UpdateAppletHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/applet/search",
		Handler: handler.SearchAppletHandler(ctx),
	})

	// Admin: 封禁用户
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/block/add",
		Handler: handler.BlockUserHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/block/del",
		Handler: handler.UnblockUserHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/block/search",
		Handler: handler.SearchBlockUserHandler(ctx),
	})

	// Admin: 重置用户密码
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/admin/user/password/reset",
		Handler: handler.AdminResetUserPasswordHandler(ctx),
	})

	// Admin: 客户端配置管理
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/client_config/set",
		Handler: handler.AdminSetClientConfigHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/client_config/del",
		Handler: handler.AdminDelClientConfigHandler(ctx),
	})

	// Admin: 统计
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistic/new_user_count",
		Handler: handler.NewUserCountHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/statistic/login_user_count",
		Handler: handler.LoginUserCountHandler(ctx),
	})

	// Admin: 版本管理
	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/application/add_version",
		Handler: handler.AddApplicationVersionHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/application/update_version",
		Handler: handler.UpdateApplicationVersionHandler(ctx),
	})

	server.AddRoute(rest.Route{
		Method:  "POST",
		Path:    "/application/delete_version",
		Handler: handler.DeleteApplicationVersionHandler(ctx),
	})

}
