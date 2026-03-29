package server

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/logic"
	"fdim/rpc/user/internal/svc"
)

type UserServer struct {
	user.UnimplementedUserServer
	svcCtx *svc.ServiceContext
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

// GetDesignateUsers 获取指定用户信息
func (s *UserServer) GetDesignateUsers(ctx context.Context, req *user.GetDesignateUsersReq) (*user.GetDesignateUsersResp, error) {
	l := logic.NewGetDesignateUsersLogic(ctx, s.svcCtx)
	return l.GetDesignateUsers(req)
}

// UpdateUserInfo 更新用户信息
func (s *UserServer) UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoReq) (*user.UpdateUserInfoResp, error) {
	l := logic.NewUpdateUserInfoLogic(ctx, s.svcCtx)
	return l.UpdateUserInfo(req)
}

// UpdateUserInfoEx 更新用户信息（扩展）
func (s *UserServer) UpdateUserInfoEx(ctx context.Context, req *user.UpdateUserInfoExReq) (*user.UpdateUserInfoExResp, error) {
	l := logic.NewUpdateUserInfoExLogic(ctx, s.svcCtx)
	return l.UpdateUserInfoEx(req)
}

// SetGlobalRecvMessageOpt 设置全局接收消息选项
func (s *UserServer) SetGlobalRecvMessageOpt(ctx context.Context, req *user.SetGlobalRecvMessageOptReq) (*user.SetGlobalRecvMessageOptResp, error) {
	l := logic.NewSetGlobalRecvMessageOptLogic(ctx, s.svcCtx)
	return l.SetGlobalRecvMessageOpt(req)
}

// AccountCheck 账户检?
func (s *UserServer) AccountCheck(ctx context.Context, req *user.AccountCheckReq) (*user.AccountCheckResp, error) {
	l := logic.NewAccountCheckLogic(ctx, s.svcCtx)
	return l.AccountCheck(req)
}

// GetPaginationUsers 分页获取用户
func (s *UserServer) GetPaginationUsers(ctx context.Context, req *user.GetPaginationUsersReq) (*user.GetPaginationUsersResp, error) {
	l := logic.NewGetPaginationUsersLogic(ctx, s.svcCtx)
	return l.GetPaginationUsers(req)
}

// SearchUserFullInfo 搜索用户完整信息
func (s *UserServer) SearchUserFullInfo(ctx context.Context, req *user.SearchUserFullInfoReq) (*user.SearchUserFullInfoResp, error) {
	l := logic.NewSearchUserFullInfoLogic(ctx, s.svcCtx)
	return l.SearchUserFullInfo(req)
}

// UserRegister 用户注册
func (s *UserServer) UserRegister(ctx context.Context, req *user.UserRegisterReq) (*user.UserRegisterResp, error) {
	l := logic.NewUserRegisterLogic(ctx, s.svcCtx)
	return l.UserRegister(req)
}

// GetGlobalRecvMessageOpt 获取全局接收消息选项
func (s *UserServer) GetGlobalRecvMessageOpt(ctx context.Context, req *user.GetGlobalRecvMessageOptReq) (*user.GetGlobalRecvMessageOptResp, error) {
	l := logic.NewGetGlobalRecvMessageOptLogic(ctx, s.svcCtx)
	return l.GetGlobalRecvMessageOpt(req)
}

// GetAllUserID 获取所有用户ID
func (s *UserServer) GetAllUserID(ctx context.Context, req *user.GetAllUserIDReq) (*user.GetAllUserIDResp, error) {
	l := logic.NewGetAllUserIDLogic(ctx, s.svcCtx)
	return l.GetAllUserID(req)
}

// GetUserStatus 获取用户状?
func (s *UserServer) GetUserStatus(ctx context.Context, req *user.GetUserStatusReq) (*user.GetUserStatusResp, error) {
	l := logic.NewGetUserStatusLogic(ctx, s.svcCtx)
	return l.GetUserStatus(req)
}

// SetUserStatus 设置用户状?
func (s *UserServer) SetUserStatus(ctx context.Context, req *user.SetUserStatusReq) (*user.SetUserStatusResp, error) {
	l := logic.NewSetUserStatusLogic(ctx, s.svcCtx)
	return l.SetUserStatus(req)
}

// SetUserOnlineStatus 设置用户在线状?
func (s *UserServer) SetUserOnlineStatus(ctx context.Context, req *user.SetUserOnlineStatusReq) (*user.SetUserOnlineStatusResp, error) {
	l := logic.NewSetUserOnlineStatusLogic(ctx, s.svcCtx)
	return l.SetUserOnlineStatus(req)
}

// GetAllOnlineUsers 获取所有在线用?
func (s *UserServer) GetAllOnlineUsers(ctx context.Context, req *user.GetAllOnlineUsersReq) (*user.GetAllOnlineUsersResp, error) {
	l := logic.NewGetAllOnlineUsersLogic(ctx, s.svcCtx)
	return l.GetAllOnlineUsers(req)
}

// UserRegisterCount 用户注册统计
func (s *UserServer) UserRegisterCount(ctx context.Context, req *user.UserRegisterCountReq) (*user.UserRegisterCountResp, error) {
	l := logic.NewUserRegisterCountLogic(ctx, s.svcCtx)
	return l.UserRegisterCount(req)
}

// SubscribeOrCancelUsersStatus 订阅或取消订阅用户状?
func (s *UserServer) SubscribeOrCancelUsersStatus(ctx context.Context, req *user.SubscribeOrCancelUsersStatusReq) (*user.SubscribeOrCancelUsersStatusResp, error) {
	l := logic.NewSubscribeOrCancelUsersStatusLogic(ctx, s.svcCtx)
	return l.SubscribeOrCancelUsersStatus(req)
}

// GetSubscribeUsersStatus 获取订阅用户状?
func (s *UserServer) GetSubscribeUsersStatus(ctx context.Context, req *user.GetSubscribeUsersStatusReq) (*user.GetSubscribeUsersStatusResp, error) {
	l := logic.NewGetSubscribeUsersStatusLogic(ctx, s.svcCtx)
	return l.GetSubscribeUsersStatus(req)
}

// ProcessUserCommandAdd 添加用户命令
func (s *UserServer) ProcessUserCommandAdd(ctx context.Context, req *user.ProcessUserCommandAddReq) (*user.ProcessUserCommandAddResp, error) {
	l := logic.NewProcessUserCommandAddLogic(ctx, s.svcCtx)
	return l.ProcessUserCommandAdd(req)
}

// ProcessUserCommandUpdate 更新用户命令
func (s *UserServer) ProcessUserCommandUpdate(ctx context.Context, req *user.ProcessUserCommandUpdateReq) (*user.ProcessUserCommandUpdateResp, error) {
	l := logic.NewProcessUserCommandUpdateLogic(ctx, s.svcCtx)
	return l.ProcessUserCommandUpdate(req)
}

// ProcessUserCommandDelete 删除用户命令
func (s *UserServer) ProcessUserCommandDelete(ctx context.Context, req *user.ProcessUserCommandDeleteReq) (*user.ProcessUserCommandDeleteResp, error) {
	l := logic.NewProcessUserCommandDeleteLogic(ctx, s.svcCtx)
	return l.ProcessUserCommandDelete(req)
}

// ProcessUserCommandGet 获取用户命令
func (s *UserServer) ProcessUserCommandGet(ctx context.Context, req *user.ProcessUserCommandGetReq) (*user.ProcessUserCommandGetResp, error) {
	l := logic.NewProcessUserCommandGetLogic(ctx, s.svcCtx)
	return l.ProcessUserCommandGet(req)
}

// ProcessUserCommandGetAll 获取所有用户命?
func (s *UserServer) ProcessUserCommandGetAll(ctx context.Context, req *user.ProcessUserCommandGetAllReq) (*user.ProcessUserCommandGetAllResp, error) {
	l := logic.NewProcessUserCommandGetAllLogic(ctx, s.svcCtx)
	return l.ProcessUserCommandGetAll(req)
}

// AddNotificationAccount 添加通知账户
func (s *UserServer) AddNotificationAccount(ctx context.Context, req *user.AddNotificationAccountReq) (*user.AddNotificationAccountResp, error) {
	l := logic.NewAddNotificationAccountLogic(ctx, s.svcCtx)
	return l.AddNotificationAccount(req)
}

// UpdateNotificationAccountInfo 更新通知账户信息
func (s *UserServer) UpdateNotificationAccountInfo(ctx context.Context, req *user.UpdateNotificationAccountInfoReq) (*user.UpdateNotificationAccountInfoResp, error) {
	l := logic.NewUpdateNotificationAccountInfoLogic(ctx, s.svcCtx)
	return l.UpdateNotificationAccountInfo(req)
}

// SearchNotificationAccount 搜索通知账户
func (s *UserServer) SearchNotificationAccount(ctx context.Context, req *user.SearchNotificationAccountReq) (*user.SearchNotificationAccountResp, error) {
	l := logic.NewSearchNotificationAccountLogic(ctx, s.svcCtx)
	return l.SearchNotificationAccount(req)
}

// GetNotificationAccount 获取通知账户
func (s *UserServer) GetNotificationAccount(ctx context.Context, req *user.GetNotificationAccountReq) (*user.GetNotificationAccountResp, error) {
	l := logic.NewGetNotificationAccountLogic(ctx, s.svcCtx)
	return l.GetNotificationAccount(req)
}

// SortQuery 排序查询
func (s *UserServer) SortQuery(ctx context.Context, req *user.SortQueryReq) (*user.SortQueryResp, error) {
	l := logic.NewSortQueryLogic(ctx, s.svcCtx)
	return l.SortQuery(req)
}

// GetUserClientConfig 获取用户客户端配?
func (s *UserServer) GetUserClientConfig(ctx context.Context, req *user.GetUserClientConfigReq) (*user.GetUserClientConfigResp, error) {
	l := logic.NewGetUserClientConfigLogic(ctx, s.svcCtx)
	return l.GetUserClientConfig(req)
}

// SetUserClientConfig 设置用户客户端配?
func (s *UserServer) SetUserClientConfig(ctx context.Context, req *user.SetUserClientConfigReq) (*user.SetUserClientConfigResp, error) {
	l := logic.NewSetUserClientConfigLogic(ctx, s.svcCtx)
	return l.SetUserClientConfig(req)
}

// DelUserClientConfig 删除用户客户端配?
func (s *UserServer) DelUserClientConfig(ctx context.Context, req *user.DelUserClientConfigReq) (*user.DelUserClientConfigResp, error) {
	l := logic.NewDelUserClientConfigLogic(ctx, s.svcCtx)
	return l.DelUserClientConfig(req)
}

// PageUserClientConfig 分页获取用户客户端配?
func (s *UserServer) PageUserClientConfig(ctx context.Context, req *user.PageUserClientConfigReq) (*user.PageUserClientConfigResp, error) {
	l := logic.NewPageUserClientConfigLogic(ctx, s.svcCtx)
	return l.PageUserClientConfig(req)
}
