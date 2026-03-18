package server

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/logic"
	"fdim/rpc/user/internal/svc"
)

type GroupServer struct {
	user.UnimplementedGroupServer
	svcCtx *svc.ServiceContext
}

func NewGroupServer(svcCtx *svc.ServiceContext) *GroupServer {
	return &GroupServer{
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建群组
func (s *GroupServer) CreateGroup(ctx context.Context, req *user.CreateGroupReq) (*user.CreateGroupResp, error) {
	l := logic.NewCreateGroupLogic(ctx, s.svcCtx)
	return l.CreateGroup(req)
}

// GetGroupsInfo 获取群组信息
func (s *GroupServer) GetGroupsInfo(ctx context.Context, req *user.GetGroupsInfoReq) (*user.GetGroupsInfoResp, error) {
	l := logic.NewGetGroupsInfoLogic(ctx, s.svcCtx)
	return l.GetGroupsInfo(req)
}

// SetGroupInfo 设置群组信息
func (s *GroupServer) SetGroupInfo(ctx context.Context, req *user.SetGroupInfoReq) (*user.SetGroupInfoResp, error) {
	l := logic.NewSetGroupInfoLogic(ctx, s.svcCtx)
	return l.SetGroupInfo(req)
}

// GetGroupApplicationList 获取群组申请列表
func (s *GroupServer) GetGroupApplicationList(ctx context.Context, req *user.GetGroupApplicationListReq) (*user.GetGroupApplicationListResp, error) {
	l := logic.NewGetGroupApplicationListLogic(ctx, s.svcCtx)
	return l.GetGroupApplicationList(req)
}

// GetUserReqApplicationList 获取用户申请列表
func (s *GroupServer) GetUserReqApplicationList(ctx context.Context, req *user.GetUserReqApplicationListReq) (*user.GetUserReqApplicationListResp, error) {
	l := logic.NewGetUserReqApplicationListLogic(ctx, s.svcCtx)
	return l.GetUserReqApplicationList(req)
}

// TransferGroupOwner 转移群主
func (s *GroupServer) TransferGroupOwner(ctx context.Context, req *user.TransferGroupOwnerReq) (*user.TransferGroupOwnerResp, error) {
	l := logic.NewTransferGroupOwnerLogic(ctx, s.svcCtx)
	return l.TransferGroupOwner(req)
}

// JoinGroup 加入群组
func (s *GroupServer) JoinGroup(ctx context.Context, req *user.JoinGroupReq) (*user.JoinGroupResp, error) {
	l := logic.NewJoinGroupLogic(ctx, s.svcCtx)
	return l.JoinGroup(req)
}

// GroupApplicationResponse 群组申请响应
func (s *GroupServer) GroupApplicationResponse(ctx context.Context, req *user.GroupApplicationResponseReq) (*user.GroupApplicationResponseResp, error) {
	l := logic.NewGroupApplicationResponseLogic(ctx, s.svcCtx)
	return l.GroupApplicationResponse(req)
}

// QuitGroup 退出群组
func (s *GroupServer) QuitGroup(ctx context.Context, req *user.QuitGroupReq) (*user.QuitGroupResp, error) {
	l := logic.NewQuitGroupLogic(ctx, s.svcCtx)
	return l.QuitGroup(req)
}

// DismissGroup 解散群组
func (s *GroupServer) DismissGroup(ctx context.Context, req *user.DismissGroupReq) (*user.DismissGroupResp, error) {
	l := logic.NewDismissGroupLogic(ctx, s.svcCtx)
	return l.DismissGroup(req)
}

// MuteGroupMember 禁言群成员
func (s *GroupServer) MuteGroupMember(ctx context.Context, req *user.MuteGroupMemberReq) (*user.MuteGroupMemberResp, error) {
	l := logic.NewMuteGroupMemberLogic(ctx, s.svcCtx)
	return l.MuteGroupMember(req)
}

// CancelMuteGroupMember 取消禁言群成员
func (s *GroupServer) CancelMuteGroupMember(ctx context.Context, req *user.CancelMuteGroupMemberReq) (*user.CancelMuteGroupMemberResp, error) {
	l := logic.NewCancelMuteGroupMemberLogic(ctx, s.svcCtx)
	return l.CancelMuteGroupMember(req)
}

// MuteGroup 禁言群组
func (s *GroupServer) MuteGroup(ctx context.Context, req *user.MuteGroupReq) (*user.MuteGroupResp, error) {
	l := logic.NewMuteGroupLogic(ctx, s.svcCtx)
	return l.MuteGroup(req)
}

// CancelMuteGroup 取消禁言群组
func (s *GroupServer) CancelMuteGroup(ctx context.Context, req *user.CancelMuteGroupReq) (*user.CancelMuteGroupResp, error) {
	l := logic.NewCancelMuteGroupLogic(ctx, s.svcCtx)
	return l.CancelMuteGroup(req)
}

// SetGroupMemberInfo 设置群成员信息
func (s *GroupServer) SetGroupMemberInfo(ctx context.Context, req *user.SetGroupMemberInfoReq) (*user.SetGroupMemberInfoResp, error) {
	l := logic.NewSetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.SetGroupMemberInfo(req)
}

// GetGroupMemberList 获取群成员列表
func (s *GroupServer) GetGroupMemberList(ctx context.Context, req *user.GetGroupMemberListReq) (*user.GetGroupMemberListResp, error) {
	l := logic.NewGetGroupMemberListLogic(ctx, s.svcCtx)
	return l.GetGroupMemberList(req)
}

// GetGroupMembersInfo 获取群成员信息
func (s *GroupServer) GetGroupMembersInfo(ctx context.Context, req *user.GetGroupMembersInfoReq) (*user.GetGroupMembersInfoResp, error) {
	l := logic.NewGetGroupMembersInfoLogic(ctx, s.svcCtx)
	return l.GetGroupMembersInfo(req)
}

// KickGroupMember 踢出群成员
func (s *GroupServer) KickGroupMember(ctx context.Context, req *user.KickGroupMemberReq) (*user.KickGroupMemberResp, error) {
	l := logic.NewKickGroupMemberLogic(ctx, s.svcCtx)
	return l.KickGroupMember(req)
}

// GetJoinedGroupList 获取已加入的群组列表
func (s *GroupServer) GetJoinedGroupList(ctx context.Context, req *user.GetJoinedGroupListReq) (*user.GetJoinedGroupListResp, error) {
	l := logic.NewGetJoinedGroupListLogic(ctx, s.svcCtx)
	return l.GetJoinedGroupList(req)
}

// InviteUserToGroup 邀请用户加入群组
func (s *GroupServer) InviteUserToGroup(ctx context.Context, req *user.InviteUserToGroupReq) (*user.InviteUserToGroupResp, error) {
	l := logic.NewInviteUserToGroupLogic(ctx, s.svcCtx)
	return l.InviteUserToGroup(req)
}

// GetGroupAllMember 获取群组所有成员
func (s *GroupServer) GetGroupAllMember(ctx context.Context, req *user.GetGroupAllMemberReq) (*user.GetGroupAllMemberResp, error) {
	l := logic.NewGetGroupAllMemberLogic(ctx, s.svcCtx)
	return l.GetGroupAllMember(req)
}

// GetGroupApplicationUnhandledCount 获取未处理的群组申请数量
func (s *GroupServer) GetGroupApplicationUnhandledCount(ctx context.Context, req *user.GetGroupApplicationUnhandledCountReq) (*user.GetGroupApplicationUnhandledCountResp, error) {
	l := logic.NewGetGroupApplicationUnhandledCountLogic(ctx, s.svcCtx)
	return l.GetGroupApplicationUnhandledCount(req)
}

// GetGroupMemberUserIDs 获取群成员用户ID列表
func (s *GroupServer) GetGroupMemberUserIDs(ctx context.Context, req *user.GetGroupMemberUserIDsReq) (*user.GetGroupMemberUserIDsResp, error) {
	l := logic.NewGetGroupMemberUserIDsLogic(ctx, s.svcCtx)
	return l.GetGroupMemberUserIDs(req)
}

// GetFullGroupMemberUserIDs 获取完整群成员用户ID列表
func (s *GroupServer) GetFullGroupMemberUserIDs(ctx context.Context, req *user.GetFullGroupMemberUserIDsReq) (*user.GetFullGroupMemberUserIDsResp, error) {
	l := logic.NewGetFullGroupMemberUserIDsLogic(ctx, s.svcCtx)
	return l.GetFullGroupMemberUserIDs(req)
}

// GetFullJoinGroupIDs 获取完整加入的群组ID列表
func (s *GroupServer) GetFullJoinGroupIDs(ctx context.Context, req *user.GetFullJoinGroupIDsReq) (*user.GetFullJoinGroupIDsResp, error) {
	l := logic.NewGetFullJoinGroupIDsLogic(ctx, s.svcCtx)
	return l.GetFullJoinGroupIDs(req)
}

// GetIncrementalGroupMember 获取增量群成员
func (s *GroupServer) GetIncrementalGroupMember(ctx context.Context, req *user.GetIncrementalGroupMemberReq) (*user.GetIncrementalGroupMemberResp, error) {
	l := logic.NewGetIncrementalGroupMemberLogic(ctx, s.svcCtx)
	return l.GetIncrementalGroupMember(req)
}

// GetIncrementalJoinGroup 获取增量加入的群组
func (s *GroupServer) GetIncrementalJoinGroup(ctx context.Context, req *user.GetIncrementalJoinGroupReq) (*user.GetIncrementalJoinGroupResp, error) {
	l := logic.NewGetIncrementalJoinGroupLogic(ctx, s.svcCtx)
	return l.GetIncrementalJoinGroup(req)
}

// GroupCreateCount 群组创建统计
func (s *GroupServer) GroupCreateCount(ctx context.Context, req *user.GroupCreateCountReq) (*user.GroupCreateCountResp, error) {
	l := logic.NewGroupCreateCountLogic(ctx, s.svcCtx)
	return l.GroupCreateCount(req)
}

// BatchGetIncrementalGroupMember 批量获取增量群成员
func (s *GroupServer) BatchGetIncrementalGroupMember(ctx context.Context, req *user.BatchGetIncrementalGroupMemberReq) (*user.BatchGetIncrementalGroupMemberResp, error) {
	l := logic.NewBatchGetIncrementalGroupMemberLogic(ctx, s.svcCtx)
	return l.BatchGetIncrementalGroupMember(req)
}

// GetGroupAbstractInfo 获取群组抽象信息
func (s *GroupServer) GetGroupAbstractInfo(ctx context.Context, req *user.GetGroupAbstractInfoReq) (*user.GetGroupAbstractInfoResp, error) {
	l := logic.NewGetGroupAbstractInfoLogic(ctx, s.svcCtx)
	return l.GetGroupAbstractInfo(req)
}

// SetGroupInfoEx 设置群组信息（扩展）
func (s *GroupServer) SetGroupInfoEx(ctx context.Context, req *user.SetGroupInfoExReq) (*user.SetGroupInfoExResp, error) {
	l := logic.NewSetGroupInfoExLogic(ctx, s.svcCtx)
	return l.SetGroupInfoEx(req)
}

// GetGroups 获取群组列表
func (s *GroupServer) GetGroups(ctx context.Context, req *user.GetGroupsReq) (*user.GetGroupsResp, error) {
	l := logic.NewGetGroupsLogic(ctx, s.svcCtx)
	return l.GetGroups(req)
}

// GetGroupMembersCMS 获取群成员CMS
func (s *GroupServer) GetGroupMembersCMS(ctx context.Context, req *user.GetGroupMembersCMSReq) (*user.GetGroupMembersCMSResp, error) {
	l := logic.NewGetGroupMembersCMSLogic(ctx, s.svcCtx)
	return l.GetGroupMembersCMS(req)
}

// GetUserInGroupMembers 获取用户在群成员中的信息
func (s *GroupServer) GetUserInGroupMembers(ctx context.Context, req *user.GetUserInGroupMembersReq) (*user.GetUserInGroupMembersResp, error) {
	l := logic.NewGetUserInGroupMembersLogic(ctx, s.svcCtx)
	return l.GetUserInGroupMembers(req)
}

// GetGroupMemberRoleLevel 获取群成员角色级别
func (s *GroupServer) GetGroupMemberRoleLevel(ctx context.Context, req *user.GetGroupMemberRoleLevelReq) (*user.GetGroupMemberRoleLevelResp, error) {
	l := logic.NewGetGroupMemberRoleLevelLogic(ctx, s.svcCtx)
	return l.GetGroupMemberRoleLevel(req)
}

// GetGroupInfoCache 获取群组信息缓存
func (s *GroupServer) GetGroupInfoCache(ctx context.Context, req *user.GetGroupInfoCacheReq) (*user.GetGroupInfoCacheResp, error) {
	l := logic.NewGetGroupInfoCacheLogic(ctx, s.svcCtx)
	return l.GetGroupInfoCache(req)
}

// GetGroupMemberCache 获取群成员缓存
func (s *GroupServer) GetGroupMemberCache(ctx context.Context, req *user.GetGroupMemberCacheReq) (*user.GetGroupMemberCacheResp, error) {
	l := logic.NewGetGroupMemberCacheLogic(ctx, s.svcCtx)
	return l.GetGroupMemberCache(req)
}

// NotificationUserInfoUpdate 通知用户信息更新
func (s *GroupServer) NotificationUserInfoUpdate(ctx context.Context, req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	l := logic.NewNotificationUserInfoUpdateGroupLogic(ctx, s.svcCtx)
	return l.NotificationUserInfoUpdate(req)
}

// GetGroupUsersReqApplicationList 获取群组用户申请列表
func (s *GroupServer) GetGroupUsersReqApplicationList(ctx context.Context, req *user.GetGroupUsersReqApplicationListReq) (*user.GetGroupUsersReqApplicationListResp, error) {
	l := logic.NewGetGroupUsersReqApplicationListLogic(ctx, s.svcCtx)
	return l.GetGroupUsersReqApplicationList(req)
}

// GetSpecifiedUserGroupRequestInfo 获取指定用户群组请求信息
func (s *GroupServer) GetSpecifiedUserGroupRequestInfo(ctx context.Context, req *user.GetSpecifiedUserGroupRequestInfoReq) (*user.GetSpecifiedUserGroupRequestInfoResp, error) {
	l := logic.NewGetSpecifiedUserGroupRequestInfoLogic(ctx, s.svcCtx)
	return l.GetSpecifiedUserGroupRequestInfo(req)
}
