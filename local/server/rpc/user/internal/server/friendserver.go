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

// ApplyToAddFriend Ӻ
func (s *FriendServer) ApplyToAddFriend(ctx context.Context, req *user.ApplyToAddFriendReq) (*user.ApplyToAddFriendResp, error) {
	l := logic.NewApplyToAddFriendLogic(ctx, s.svcCtx)
	return l.ApplyToAddFriend(req)
}

// ImportFriends 
func (s *FriendServer) ImportFriends(ctx context.Context, req *user.ImportFriendReq) (*user.ImportFriendResp, error) {
	l := logic.NewImportFriendsLogic(ctx, s.svcCtx)
	return l.ImportFriends(req)
}

// RespondFriendApply Ӧ
func (s *FriendServer) RespondFriendApply(ctx context.Context, req *user.RespondFriendApplyReq) (*user.RespondFriendApplyResp, error) {
	l := logic.NewRespondFriendApplyLogic(ctx, s.svcCtx)
	return l.RespondFriendApply(req)
}

// DeleteFriend ɾ
func (s *FriendServer) DeleteFriend(ctx context.Context, req *user.DeleteFriendReq) (*user.DeleteFriendResp, error) {
	l := logic.NewDeleteFriendLogic(ctx, s.svcCtx)
	return l.DeleteFriend(req)
}

// GetPaginationFriends ҳȡ
func (s *FriendServer) GetPaginationFriends(ctx context.Context, req *user.GetPaginationFriendsReq) (*user.GetPaginationFriendsResp, error) {
	l := logic.NewGetPaginationFriendsLogic(ctx, s.svcCtx)
	return l.GetPaginationFriends(req)
}

// GetFriendIDs ȡIDб
func (s *FriendServer) GetFriendIDs(ctx context.Context, req *user.GetFriendIDsReq) (*user.GetFriendIDsResp, error) {
	l := logic.NewGetFriendIDsLogic(ctx, s.svcCtx)
	return l.GetFriendIDs(req)
}

// GetSpecifiedFriends ȡָϢ
func (s *FriendServer) GetSpecifiedFriends(ctx context.Context, req *user.GetSpecifiedFriendsInfoReq) (*user.GetSpecifiedFriendsInfoResp, error) {
	l := logic.NewGetSpecifiedFriendsLogic(ctx, s.svcCtx)
	return l.GetSpecifiedFriends(req)
}

// GetPaginationFriendsApplyTo ȡбյģ
func (s *FriendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *user.GetPaginationFriendsApplyToReq) (*user.GetPaginationFriendsApplyToResp, error) {
	l := logic.NewGetFriendApplyListLogic(ctx, s.svcCtx)
	return l.GetFriendApplyList(req)
}

// GetPaginationFriendsApplyFrom ȡԼб
func (s *FriendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *user.GetPaginationFriendsApplyFromReq) (*user.GetPaginationFriendsApplyFromResp, error) {
	l := logic.NewGetSelfApplyListLogic(ctx, s.svcCtx)
	return l.GetSelfApplyList(req)
}

// IsFriend ǷΪ
func (s *FriendServer) IsFriend(ctx context.Context, req *user.IsFriendReq) (*user.IsFriendResp, error) {
	l := logic.NewIsFriendLogic(ctx, s.svcCtx)
	return l.IsFriend(req)
}

// GetFriendApplyUnhandledCount ȡδĺ
func (s *FriendServer) GetFriendApplyUnhandledCount(ctx context.Context, req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	l := logic.NewGetFriendApplyUnhandledCountLogic(ctx, s.svcCtx)
	return l.GetFriendApplyUnhandledCount(req)
}

// GetSelfUnhandledApplyCount ȡԼδ
func (s *FriendServer) GetSelfUnhandledApplyCount(ctx context.Context, req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	l := logic.NewGetSelfUnhandledApplyCountLogic(ctx, s.svcCtx)
	return l.GetSelfUnhandledApplyCount(req)
}

// SetFriendRemark úѱע
func (s *FriendServer) SetFriendRemark(ctx context.Context, req *user.SetFriendRemarkReq) (*user.SetFriendRemarkResp, error) {
	l := logic.NewSetFriendRemarkLogic(ctx, s.svcCtx)
	return l.SetFriendRemark(req)
}

// UpdateFriends ºϢ
func (s *FriendServer) UpdateFriends(ctx context.Context, req *user.UpdateFriendsReq) (*user.UpdateFriendsResp, error) {
	l := logic.NewUpdateFriendsLogic(ctx, s.svcCtx)
	return l.UpdateFriends(req)
}

// AddBlack Ӻ
func (s *FriendServer) AddBlack(ctx context.Context, req *user.AddBlackReq) (*user.AddBlackResp, error) {
	l := logic.NewAddBlackLogic(ctx, s.svcCtx)
	return l.AddBlack(req)
}

// RemoveBlack Ƴ
func (s *FriendServer) RemoveBlack(ctx context.Context, req *user.RemoveBlackReq) (*user.RemoveBlackResp, error) {
	l := logic.NewRemoveBlackLogic(ctx, s.svcCtx)
	return l.RemoveBlack(req)
}

// GetPaginationBlacks ҳȡ
func (s *FriendServer) GetPaginationBlacks(ctx context.Context, req *user.GetPaginationBlacksReq) (*user.GetPaginationBlacksResp, error) {
	l := logic.NewGetPaginationBlacksLogic(ctx, s.svcCtx)
	return l.GetPaginationBlacks(req)
}

// IsBlack Ƿں
func (s *FriendServer) IsBlack(ctx context.Context, req *user.IsBlackReq) (*user.IsBlackResp, error) {
	l := logic.NewIsBlackLogic(ctx, s.svcCtx)
	return l.IsBlack(req)
}

// GetSpecifiedBlacks ȡָ
func (s *FriendServer) GetSpecifiedBlacks(ctx context.Context, req *user.GetSpecifiedBlacksReq) (*user.GetSpecifiedBlacksResp, error) {
	l := logic.NewGetSpecifiedBlacksLogic(ctx, s.svcCtx)
	return l.GetSpecifiedBlacks(req)
}

// GetFullFriendUserIDs ȡûIDб
func (s *FriendServer) GetFullFriendUserIDs(ctx context.Context, req *user.GetFullFriendUserIDsReq) (*user.GetFullFriendUserIDsResp, error) {
	l := logic.NewGetFullFriendUserIDsLogic(ctx, s.svcCtx)
	return l.GetFullFriendUserIDs(req)
}

// GetIncrementalFriends ȡ
func (s *FriendServer) GetIncrementalFriends(ctx context.Context, req *user.GetIncrementalFriendsReq) (*user.GetIncrementalFriendsResp, error) {
	l := logic.NewGetIncrementalFriendsLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriends(req)
}

// GetDesignatedFriends ȡָϢ
func (s *FriendServer) GetDesignatedFriends(ctx context.Context, req *user.GetDesignatedFriendsReq) (*user.GetDesignatedFriendsResp, error) {
	l := logic.NewGetDesignatedFriendsLogic(ctx, s.svcCtx)
	return l.GetDesignatedFriends(req)
}

// GetDesignatedFriendsApply ȡָ
func (s *FriendServer) GetDesignatedFriendsApply(ctx context.Context, req *user.GetDesignatedFriendsApplyReq) (*user.GetDesignatedFriendsApplyResp, error) {
	l := logic.NewGetDesignatedFriendsApplyLogic(ctx, s.svcCtx)
	return l.GetDesignatedFriendsApply(req)
}

// GetFriendInfo ȡϢ
func (s *FriendServer) GetFriendInfo(ctx context.Context, req *user.GetFriendInfoReq) (*user.GetFriendInfoResp, error) {
	l := logic.NewGetFriendInfoLogic(ctx, s.svcCtx)
	return l.GetFriendInfo(req)
}

// GetIncrementalBlacks ȡ
func (s *FriendServer) GetIncrementalBlacks(ctx context.Context, req *user.GetIncrementalBlacksReq) (*user.GetIncrementalBlacksResp, error) {
	l := logic.NewGetIncrementalBlacksLogic(ctx, s.svcCtx)
	return l.GetIncrementalBlacks(req)
}

// GetIncrementalFriendsApplyTo ȡ루յ
func (s *FriendServer) GetIncrementalFriendsApplyTo(ctx context.Context, req *user.GetIncrementalFriendsApplyToReq) (*user.GetIncrementalFriendsApplyToResp, error) {
	l := logic.NewGetIncrementalFriendsApplyToLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriendsApplyTo(req)
}

// GetIncrementalFriendsApplyFrom ȡ루
func (s *FriendServer) GetIncrementalFriendsApplyFrom(ctx context.Context, req *user.GetIncrementalFriendsApplyFromReq) (*user.GetIncrementalFriendsApplyFromResp, error) {
	l := logic.NewGetIncrementalFriendsApplyFromLogic(ctx, s.svcCtx)
	return l.GetIncrementalFriendsApplyFrom(req)
}

// NotificationUserInfoUpdate ֪ͨûϢ
func (s *FriendServer) NotificationUserInfoUpdate(ctx context.Context, req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	l := logic.NewNotificationUserInfoUpdateLogic(ctx, s.svcCtx)
	return l.NotificationUserInfoUpdate(req)
}
