package server

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/logic"
	"fdim/rpc/user/internal/svc"
)

type FriendServer struct {
	user.UnimplementedFriendServer
	svcCtx *svc.ServiceContext
}

func NewFriendServer(svcCtx *svc.ServiceContext) *FriendServer {
	return &FriendServer{
		svcCtx: svcCtx,
	}
}

// ApplyToAddFriend 申请添加好友
func (s *FriendServer) ApplyToAddFriend(ctx context.Context, req *user.ApplyToAddFriendReq) (*user.ApplyToAddFriendResp, error) {
	l := logic.NewApplyToAddFriendLogic(ctx, s.svcCtx)
	return l.ApplyToAddFriend(req)
}

// ImportFriends 导入好友
func (s *FriendServer) ImportFriends(ctx context.Context, req *user.ImportFriendReq) (*user.ImportFriendResp, error) {
	l := logic.NewImportFriendsLogic(ctx, s.svcCtx)
	return l.ImportFriends(req)
}

// RespondFriendApply 响应好友申请
func (s *FriendServer) RespondFriendApply(ctx context.Context, req *user.RespondFriendApplyReq) (*user.RespondFriendApplyResp, error) {
	l := logic.NewRespondFriendApplyLogic(ctx, s.svcCtx)
	return l.RespondFriendApply(req)
}

// DeleteFriend 删除好友
func (s *FriendServer) DeleteFriend(ctx context.Context, req *user.DeleteFriendReq) (*user.DeleteFriendResp, error) {
	l := logic.NewDeleteFriendLogic(ctx, s.svcCtx)
	return l.DeleteFriend(req)
}

// GetPaginationFriends 分页获取好友
func (s *FriendServer) GetPaginationFriends(ctx context.Context, req *user.GetPaginationFriendsReq) (*user.GetPaginationFriendsResp, error) {
	l := logic.NewGetPaginationFriendsLogic(ctx, s.svcCtx)
	return l.GetPaginationFriends(req)
}

// GetFriendIDs 获取好友ID列表
func (s *FriendServer) GetFriendIDs(ctx context.Context, req *user.GetFriendIDsReq) (*user.GetFriendIDsResp, error) {
	l := logic.NewGetFriendIDsLogic(ctx, s.svcCtx)
	return l.GetFriendIDs(req)
}

// GetSpecifiedFriends 获取指定好友信息
func (s *FriendServer) GetSpecifiedFriends(ctx context.Context, req *user.GetSpecifiedFriendsInfoReq) (*user.GetSpecifiedFriendsInfoResp, error) {
	l := logic.NewGetSpecifiedFriendsLogic(ctx, s.svcCtx)
	return l.GetSpecifiedFriends(req)
}

// GetPaginationFriendsApplyTo 获取好友申请列表（收到的）
func (s *FriendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *user.GetPaginationFriendsApplyToReq) (*user.GetPaginationFriendsApplyToResp, error) {
	l := logic.NewGetFriendApplyListLogic(ctx, s.svcCtx)
	return l.GetFriendApplyList(req)
}

// GetPaginationFriendsApplyFrom 获取自己发出的申请列表
func (s *FriendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *user.GetPaginationFriendsApplyFromReq) (*user.GetPaginationFriendsApplyFromResp, error) {
	l := logic.NewGetSelfApplyListLogic(ctx, s.svcCtx)
	return l.GetSelfApplyList(req)
}

// IsFriend 检查是否为好友
func (s *FriendServer) IsFriend(ctx context.Context, req *user.IsFriendReq) (*user.IsFriendResp, error) {
	l := logic.NewIsFriendLogic(ctx, s.svcCtx)
	return l.IsFriend(req)
}

// GetFriendApplyUnhandledCount 获取未处理的好友申请数量
func (s *FriendServer) GetFriendApplyUnhandledCount(ctx context.Context, req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	l := logic.NewGetFriendApplyUnhandledCountLogic(ctx, s.svcCtx)
	return l.GetFriendApplyUnhandledCount(req)
}

// GetSelfUnhandledApplyCount 获取自己未处理的申请数量
func (s *FriendServer) GetSelfUnhandledApplyCount(ctx context.Context, req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	l := logic.NewGetSelfUnhandledApplyCountLogic(ctx, s.svcCtx)
	return l.GetSelfUnhandledApplyCount(req)
}

// SetFriendRemark 设置好友备注
func (s *FriendServer) SetFriendRemark(ctx context.Context, req *user.SetFriendRemarkReq) (*user.SetFriendRemarkResp, error) {
	l := logic.NewSetFriendRemarkLogic(ctx, s.svcCtx)
	return l.SetFriendRemark(req)
}

// UpdateFriends 更新好友信息
func (s *FriendServer) UpdateFriends(ctx context.Context, req *user.UpdateFriendsReq) (*user.UpdateFriendsResp, error) {
	l := logic.NewUpdateFriendsLogic(ctx, s.svcCtx)
	return l.UpdateFriends(req)
}

// AddBlack 添加黑名单
func (s *FriendServer) AddBlack(ctx context.Context, req *user.AddBlackReq) (*user.AddBlackResp, error) {
	l := logic.NewAddBlackLogic(ctx, s.svcCtx)
	return l.AddBlack(req)
}

// RemoveBlack 移除黑名单
func (s *FriendServer) RemoveBlack(ctx context.Context, req *user.RemoveBlackReq) (*user.RemoveBlackResp, error) {
	l := logic.NewRemoveBlackLogic(ctx, s.svcCtx)
	return l.RemoveBlack(req)
}

// GetPaginationBlacks 分页获取黑名单
func (s *FriendServer) GetPaginationBlacks(ctx context.Context, req *user.GetPaginationBlacksReq) (*user.GetPaginationBlacksResp, error) {
	l := logic.NewGetPaginationBlacksLogic(ctx, s.svcCtx)
	return l.GetPaginationBlacks(req)
}

// IsBlack 检查是否在黑名单
func (s *FriendServer) IsBlack(ctx context.Context, req *user.IsBlackReq) (*user.IsBlackResp, error) {
	l := logic.NewIsBlackLogic(ctx, s.svcCtx)
	return l.IsBlack(req)
}

// GetSpecifiedBlacks 获取指定黑名单
func (s *FriendServer) GetSpecifiedBlacks(ctx context.Context, req *user.GetSpecifiedBlacksReq) (*user.GetSpecifiedBlacksResp, error) {
	l := logic.NewGetSpecifiedBlacksLogic(ctx, s.svcCtx)
	return l.GetSpecifiedBlacks(req)
}

// GetFullFriendUserIDs 获取完整好友用户ID列表
func (s *FriendServer) GetFullFriendUserIDs(ctx context.Context, req *user.GetFullFriendUserIDsReq) (*user.GetFullFriendUserIDsResp, error) {
	l := logic.NewGetFullFriendUserIDsLogic(ctx, s.svcCtx)
	return l.GetFullFriendUserIDs(req)
}

// GetIncrementalFriends 获取增量好友
func (s *FriendServer) GetIncrementalFriends(ctx context.Context, req *user.GetIncrementalFriendsReq) (*user.GetIncrementalFriendsResp, error) {
	l := logic.NewGetIncrementalFriendsLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriends(req)
}

// GetDesignatedFriends 获取指定好友信息
func (s *FriendServer) GetDesignatedFriends(ctx context.Context, req *user.GetDesignatedFriendsReq) (*user.GetDesignatedFriendsResp, error) {
	l := logic.NewGetDesignatedFriendsLogic(ctx, s.svcCtx)
	return l.GetDesignatedFriends(req)
}

// GetDesignatedFriendsApply 获取指定好友申请
func (s *FriendServer) GetDesignatedFriendsApply(ctx context.Context, req *user.GetDesignatedFriendsApplyReq) (*user.GetDesignatedFriendsApplyResp, error) {
	l := logic.NewGetDesignatedFriendsApplyLogic(ctx, s.svcCtx)
	return l.GetDesignatedFriendsApply(req)
}

// GetFriendInfo 获取好友信息
func (s *FriendServer) GetFriendInfo(ctx context.Context, req *user.GetFriendInfoReq) (*user.GetFriendInfoResp, error) {
	l := logic.NewGetFriendInfoLogic(ctx, s.svcCtx)
	return l.GetFriendInfo(req)
}

// GetIncrementalBlacks 获取增量黑名单
func (s *FriendServer) GetIncrementalBlacks(ctx context.Context, req *user.GetIncrementalBlacksReq) (*user.GetIncrementalBlacksResp, error) {
	l := logic.NewGetIncrementalBlacksLogic(ctx, s.svcCtx)
	return l.GetIncrementalBlacks(req)
}

// GetIncrementalFriendsApplyTo 获取增量好友申请（收到）
func (s *FriendServer) GetIncrementalFriendsApplyTo(ctx context.Context, req *user.GetIncrementalFriendsApplyToReq) (*user.GetIncrementalFriendsApplyToResp, error) {
	l := logic.NewGetIncrementalFriendsApplyToLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriendsApplyTo(req)
}

// GetIncrementalFriendsApplyFrom 获取增量好友申请（发出）
func (s *FriendServer) GetIncrementalFriendsApplyFrom(ctx context.Context, req *user.GetIncrementalFriendsApplyFromReq) (*user.GetIncrementalFriendsApplyFromResp, error) {
	l := logic.NewGetIncrementalFriendsApplyFromLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriendsApplyFrom(req)
}

// NotificationUserInfoUpdate 通知用户信息更新
func (s *FriendServer) NotificationUserInfoUpdate(ctx context.Context, req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	l := logic.NewNotificationUserInfoUpdateLogic(ctx, s.svcCtx)
	return l.NotificationUserInfoUpdate(req)
}
