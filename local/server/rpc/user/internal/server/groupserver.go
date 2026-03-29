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

// CreateGroup Ⱥ
func (s *GroupServer) CreateGroup(ctx context.Context, req *user.CreateGroupReq) (*user.CreateGroupResp, error) {
	l := logic.NewCreateGroupLogic(ctx, s.svcCtx)
	return l.CreateGroup(req)
}

// GetGroupsInfo ȡȺϢ
func (s *GroupServer) GetGroupsInfo(ctx context.Context, req *user.GetGroupsInfoReq) (*user.GetGroupsInfoResp, error) {
	l := logic.NewGetGroupsInfoLogic(ctx, s.svcCtx)
	return l.GetGroupsInfo(req)
}

// SetGroupInfo ȺϢ
func (s *GroupServer) SetGroupInfo(ctx context.Context, req *user.SetGroupInfoReq) (*user.SetGroupInfoResp, error) {
	l := logic.NewSetGroupInfoLogic(ctx, s.svcCtx)
	return l.SetGroupInfo(req)
}

// GetGroupApplicationList ȡȺб
func (s *GroupServer) GetGroupApplicationList(ctx context.Context, req *user.GetGroupApplicationListReq) (*user.GetGroupApplicationListResp, error) {
	l := logic.NewGetGroupApplicationListLogic(ctx, s.svcCtx)
	return l.GetGroupApplicationList(req)
}

// GetUserReqApplicationList ȡûб
func (s *GroupServer) GetUserReqApplicationList(ctx context.Context, req *user.GetUserReqApplicationListReq) (*user.GetUserReqApplicationListResp, error) {
	l := logic.NewGetUserReqApplicationListLogic(ctx, s.svcCtx)
	return l.GetUserReqApplicationList(req)
}

// TransferGroupOwner תȺ
func (s *GroupServer) TransferGroupOwner(ctx context.Context, req *user.TransferGroupOwnerReq) (*user.TransferGroupOwnerResp, error) {
	l := logic.NewTransferGroupOwnerLogic(ctx, s.svcCtx)
	return l.TransferGroupOwner(req)
}

// JoinGroup Ⱥ
func (s *GroupServer) JoinGroup(ctx context.Context, req *user.JoinGroupReq) (*user.JoinGroupResp, error) {
	l := logic.NewJoinGroupLogic(ctx, s.svcCtx)
	return l.JoinGroup(req)
}

// GroupApplicationResponse ȺӦ
func (s *GroupServer) GroupApplicationResponse(ctx context.Context, req *user.GroupApplicationResponseReq) (*user.GroupApplicationResponseResp, error) {
	l := logic.NewGroupApplicationResponseLogic(ctx, s.svcCtx)
	return l.GroupApplicationResponse(req)
}

// QuitGroup ˳Ⱥ
func (s *GroupServer) QuitGroup(ctx context.Context, req *user.QuitGroupReq) (*user.QuitGroupResp, error) {
	l := logic.NewQuitGroupLogic(ctx, s.svcCtx)
	return l.QuitGroup(req)
}

// DismissGroup ɢȺ
func (s *GroupServer) DismissGroup(ctx context.Context, req *user.DismissGroupReq) (*user.DismissGroupResp, error) {
	l := logic.NewDismissGroupLogic(ctx, s.svcCtx)
	return l.DismissGroup(req)
}

// MuteGroupMember ȺԱ
func (s *GroupServer) MuteGroupMember(ctx context.Context, req *user.MuteGroupMemberReq) (*user.MuteGroupMemberResp, error) {
	l := logic.NewMuteGroupMemberLogic(ctx, s.svcCtx)
	return l.MuteGroupMember(req)
}

// CancelMuteGroupMember ȡȺԱ
func (s *GroupServer) CancelMuteGroupMember(ctx context.Context, req *user.CancelMuteGroupMemberReq) (*user.CancelMuteGroupMemberResp, error) {
	l := logic.NewCancelMuteGroupMemberLogic(ctx, s.svcCtx)
	return l.CancelMuteGroupMember(req)
}

// MuteGroup Ⱥ
func (s *GroupServer) MuteGroup(ctx context.Context, req *user.MuteGroupReq) (*user.MuteGroupResp, error) {
	l := logic.NewMuteGroupLogic(ctx, s.svcCtx)
	return l.MuteGroup(req)
}

// CancelMuteGroup ȡȺ
func (s *GroupServer) CancelMuteGroup(ctx context.Context, req *user.CancelMuteGroupReq) (*user.CancelMuteGroupResp, error) {
	l := logic.NewCancelMuteGroupLogic(ctx, s.svcCtx)
	return l.CancelMuteGroup(req)
}

// SetGroupMemberInfo ȺԱϢ
func (s *GroupServer) SetGroupMemberInfo(ctx context.Context, req *user.SetGroupMemberInfoReq) (*user.SetGroupMemberInfoResp, error) {
	l := logic.NewSetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.SetGroupMemberInfo(req)
}

// GetGroupMemberList ȡȺԱб
func (s *GroupServer) GetGroupMemberList(ctx context.Context, req *user.GetGroupMemberListReq) (*user.GetGroupMemberListResp, error) {
	l := logic.NewGetGroupMemberListLogic(ctx, s.svcCtx)
	return l.GetGroupMemberList(req)
}

// GetGroupMembersInfo ȡȺԱϢ
func (s *GroupServer) GetGroupMembersInfo(ctx context.Context, req *user.GetGroupMembersInfoReq) (*user.GetGroupMembersInfoResp, error) {
	l := logic.NewGetGroupMembersInfoLogic(ctx, s.svcCtx)
	return l.GetGroupMembersInfo(req)
}

// KickGroupMember ߳ȺԱ
func (s *GroupServer) KickGroupMember(ctx context.Context, req *user.KickGroupMemberReq) (*user.KickGroupMemberResp, error) {
	l := logic.NewKickGroupMemberLogic(ctx, s.svcCtx)
	return l.KickGroupMember(req)
}

// GetJoinedGroupList ȡѼȺб
func (s *GroupServer) GetJoinedGroupList(ctx context.Context, req *user.GetJoinedGroupListReq) (*user.GetJoinedGroupListResp, error) {
	l := logic.NewGetJoinedGroupListLogic(ctx, s.svcCtx)
	return l.GetJoinedGroupList(req)
}

// InviteUserToGroup ûȺ
func (s *GroupServer) InviteUserToGroup(ctx context.Context, req *user.InviteUserToGroupReq) (*user.InviteUserToGroupResp, error) {
	l := logic.NewInviteUserToGroupLogic(ctx, s.svcCtx)
	return l.InviteUserToGroup(req)
}

// GetGroupAllMember ȡȺгԱ
func (s *GroupServer) GetGroupAllMember(ctx context.Context, req *user.GetGroupAllMemberReq) (*user.GetGroupAllMemberResp, error) {
	l := logic.NewGetGroupAllMemberLogic(ctx, s.svcCtx)
	return l.GetGroupAllMember(req)
}

// GetGroupApplicationUnhandledCount ȡδȺ
func (s *GroupServer) GetGroupApplicationUnhandledCount(ctx context.Context, req *user.GetGroupApplicationUnhandledCountReq) (*user.GetGroupApplicationUnhandledCountResp, error) {
	l := logic.NewGetGroupApplicationUnhandledCountLogic(ctx, s.svcCtx)
	return l.GetGroupApplicationUnhandledCount(req)
}

// GetGroupMemberUserIDs ȡȺԱûIDб
func (s *GroupServer) GetGroupMemberUserIDs(ctx context.Context, req *user.GetGroupMemberUserIDsReq) (*user.GetGroupMemberUserIDsResp, error) {
	l := logic.NewGetGroupMemberUserIDsLogic(ctx, s.svcCtx)
	return l.GetGroupMemberUserIDs(req)
}

// GetFullGroupMemberUserIDs ȡȺԱûIDб
func (s *GroupServer) GetFullGroupMemberUserIDs(ctx context.Context, req *user.GetFullGroupMemberUserIDsReq) (*user.GetFullGroupMemberUserIDsResp, error) {
	l := logic.NewGetFullGroupMemberUserIDsLogic(ctx, s.svcCtx)
	return l.GetFullGroupMemberUserIDs(req)
}

// GetFullJoinGroupIDs ȡȺIDб
func (s *GroupServer) GetFullJoinGroupIDs(ctx context.Context, req *user.GetFullJoinGroupIDsReq) (*user.GetFullJoinGroupIDsResp, error) {
	l := logic.NewGetFullJoinGroupIDsLogic(ctx, s.svcCtx)
	return l.GetFullJoinGroupIDs(req)
}

// GetIncrementalGroupMember ȡȺԱ
func (s *GroupServer) GetIncrementalGroupMember(ctx context.Context, req *user.GetIncrementalGroupMemberReq) (*user.GetIncrementalGroupMemberResp, error) {
	l := logic.NewGetIncrementalGroupMemberLogic(ctx, s.svcCtx)
	return l.GetIncrementalGroupMember(req)
}

// GetIncrementalJoinGroup ȡȺ
func (s *GroupServer) GetIncrementalJoinGroup(ctx context.Context, req *user.GetIncrementalJoinGroupReq) (*user.GetIncrementalJoinGroupResp, error) {
	l := logic.NewGetIncrementalJoinGroupLogic(ctx, s.svcCtx)
	return l.GetIncrementalJoinGroup(req)
}

// GroupCreateCount Ⱥ鴴ͳ
func (s *GroupServer) GroupCreateCount(ctx context.Context, req *user.GroupCreateCountReq) (*user.GroupCreateCountResp, error) {
	l := logic.NewGroupCreateCountLogic(ctx, s.svcCtx)
	return l.GroupCreateCount(req)
}

// BatchGetIncrementalGroupMember ȡȺԱ
func (s *GroupServer) BatchGetIncrementalGroupMember(ctx context.Context, req *user.BatchGetIncrementalGroupMemberReq) (*user.BatchGetIncrementalGroupMemberResp, error) {
	l := logic.NewBatchGetIncrementalGroupMemberLogic(ctx, s.svcCtx)
	return l.BatchGetIncrementalGroupMember(req)
}

// GetGroupAbstractInfo ȡȺϢ
func (s *GroupServer) GetGroupAbstractInfo(ctx context.Context, req *user.GetGroupAbstractInfoReq) (*user.GetGroupAbstractInfoResp, error) {
	l := logic.NewGetGroupAbstractInfoLogic(ctx, s.svcCtx)
	return l.GetGroupAbstractInfo(req)
}

// SetGroupInfoEx ȺϢչ
func (s *GroupServer) SetGroupInfoEx(ctx context.Context, req *user.SetGroupInfoExReq) (*user.SetGroupInfoExResp, error) {
	l := logic.NewSetGroupInfoExLogic(ctx, s.svcCtx)
	return l.SetGroupInfoEx(req)
}

// GetGroups ȡȺб
func (s *GroupServer) GetGroups(ctx context.Context, req *user.GetGroupsReq) (*user.GetGroupsResp, error) {
	l := logic.NewGetGroupsLogic(ctx, s.svcCtx)
	return l.GetGroups(req)
}

// GetGroupMembersCMS ȡȺԱCMS
func (s *GroupServer) GetGroupMembersCMS(ctx context.Context, req *user.GetGroupMembersCMSReq) (*user.GetGroupMembersCMSResp, error) {
	l := logic.NewGetGroupMembersCMSLogic(ctx, s.svcCtx)
	return l.GetGroupMembersCMS(req)
}

// GetUserInGroupMembers ȡûȺԱеϢ
func (s *GroupServer) GetUserInGroupMembers(ctx context.Context, req *user.GetUserInGroupMembersReq) (*user.GetUserInGroupMembersResp, error) {
	l := logic.NewGetUserInGroupMembersLogic(ctx, s.svcCtx)
	return l.GetUserInGroupMembers(req)
}

// GetGroupMemberRoleLevel ȡȺԱɫ
func (s *GroupServer) GetGroupMemberRoleLevel(ctx context.Context, req *user.GetGroupMemberRoleLevelReq) (*user.GetGroupMemberRoleLevelResp, error) {
	l := logic.NewGetGroupMemberRoleLevelLogic(ctx, s.svcCtx)
	return l.GetGroupMemberRoleLevel(req)
}

// GetGroupInfoCache ȡȺϢ
func (s *GroupServer) GetGroupInfoCache(ctx context.Context, req *user.GetGroupInfoCacheReq) (*user.GetGroupInfoCacheResp, error) {
	l := logic.NewGetGroupInfoCacheLogic(ctx, s.svcCtx)
	return l.GetGroupInfoCache(req)
}

// GetGroupMemberCache ȡȺԱ
func (s *GroupServer) GetGroupMemberCache(ctx context.Context, req *user.GetGroupMemberCacheReq) (*user.GetGroupMemberCacheResp, error) {
	l := logic.NewGetGroupMemberCacheLogic(ctx, s.svcCtx)
	return l.GetGroupMemberCache(req)
}

// NotificationUserInfoUpdate ֪ͨûϢ
func (s *GroupServer) NotificationUserInfoUpdate(ctx context.Context, req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	l := logic.NewNotificationUserInfoUpdateGroupLogic(ctx, s.svcCtx)
	return l.NotificationUserInfoUpdate(req)
}

// GetGroupUsersReqApplicationList ȡȺûб
func (s *GroupServer) GetGroupUsersReqApplicationList(ctx context.Context, req *user.GetGroupUsersReqApplicationListReq) (*user.GetGroupUsersReqApplicationListResp, error) {
	l := logic.NewGetGroupUsersReqApplicationListLogic(ctx, s.svcCtx)
	return l.GetGroupUsersReqApplicationList(req)
}

// GetSpecifiedUserGroupRequestInfo ȡָûȺϢ
func (s *GroupServer) GetSpecifiedUserGroupRequestInfo(ctx context.Context, req *user.GetSpecifiedUserGroupRequestInfoReq) (*user.GetSpecifiedUserGroupRequestInfoResp, error) {
	l := logic.NewGetSpecifiedUserGroupRequestInfoLogic(ctx, s.svcCtx)
	return l.GetSpecifiedUserGroupRequestInfo(req)
}
