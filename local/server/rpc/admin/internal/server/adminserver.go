package server

import (
	"context"

	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/logic"
	"fdim/rpc/admin/internal/svc"
)

type AdminServer struct {
	admin.UnimplementedAdminServer
	svcCtx *svc.ServiceContext
}

func NewAdminServer(svcCtx *svc.ServiceContext) *AdminServer {
	return &AdminServer{
		svcCtx: svcCtx,
	}
}

// Login 管理员登录
func (s *AdminServer) Login(ctx context.Context, req *admin.LoginReq) (*admin.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(req)
}

// ChangePassword 修改密码
func (s *AdminServer) ChangePassword(ctx context.Context, req *admin.ChangePasswordReq) (*admin.ChangePasswordResp, error) {
	l := logic.NewChangePasswordLogic(ctx, s.svcCtx)
	return l.ChangePassword(req)
}

// AdminUpdateInfo 更新管理员信息
func (s *AdminServer) AdminUpdateInfo(ctx context.Context, req *admin.AdminUpdateInfoReq) (*admin.AdminUpdateInfoResp, error) {
	l := logic.NewAdminUpdateInfoLogic(ctx, s.svcCtx)
	return l.AdminUpdateInfo(req)
}

// GetAdminInfo 获取管理员信息
func (s *AdminServer) GetAdminInfo(ctx context.Context, req *admin.GetAdminInfoReq) (*admin.GetAdminInfoResp, error) {
	l := logic.NewGetAdminInfoLogic(ctx, s.svcCtx)
	return l.GetAdminInfo(req)
}

// AddAdminAccount 添加管理员账户
func (s *AdminServer) AddAdminAccount(ctx context.Context, req *admin.AddAdminAccountReq) (*admin.AddAdminAccountResp, error) {
	l := logic.NewAddAdminAccountLogic(ctx, s.svcCtx)
	return l.AddAdminAccount(req)
}

// ChangeAdminPassword 修改管理员密码
func (s *AdminServer) ChangeAdminPassword(ctx context.Context, req *admin.ChangeAdminPasswordReq) (*admin.ChangeAdminPasswordResp, error) {
	l := logic.NewChangeAdminPasswordLogic(ctx, s.svcCtx)
	return l.ChangeAdminPassword(req)
}

// DelAdminAccount 删除管理员账户
func (s *AdminServer) DelAdminAccount(ctx context.Context, req *admin.DelAdminAccountReq) (*admin.DelAdminAccountResp, error) {
	l := logic.NewDelAdminAccountLogic(ctx, s.svcCtx)
	return l.DelAdminAccount(req)
}

// SearchAdminAccount 搜索管理员账户
func (s *AdminServer) SearchAdminAccount(ctx context.Context, req *admin.SearchAdminAccountReq) (*admin.SearchAdminAccountResp, error) {
	l := logic.NewSearchAdminAccountLogic(ctx, s.svcCtx)
	return l.SearchAdminAccount(req)
}

// AddDefaultFriend 添加默认好友
func (s *AdminServer) AddDefaultFriend(ctx context.Context, req *admin.AddDefaultFriendReq) (*admin.AddDefaultFriendResp, error) {
	l := logic.NewAddDefaultFriendLogic(ctx, s.svcCtx)
	return l.AddDefaultFriend(req)
}

// DelDefaultFriend 删除默认好友
func (s *AdminServer) DelDefaultFriend(ctx context.Context, req *admin.DelDefaultFriendReq) (*admin.DelDefaultFriendResp, error) {
	l := logic.NewDelDefaultFriendLogic(ctx, s.svcCtx)
	return l.DelDefaultFriend(req)
}

// FindDefaultFriend 查找默认好友
func (s *AdminServer) FindDefaultFriend(ctx context.Context, req *admin.FindDefaultFriendReq) (*admin.FindDefaultFriendResp, error) {
	l := logic.NewFindDefaultFriendLogic(ctx, s.svcCtx)
	return l.FindDefaultFriend(req)
}

// SearchDefaultFriend 搜索默认好友
func (s *AdminServer) SearchDefaultFriend(ctx context.Context, req *admin.SearchDefaultFriendReq) (*admin.SearchDefaultFriendResp, error) {
	l := logic.NewSearchDefaultFriendLogic(ctx, s.svcCtx)
	return l.SearchDefaultFriend(req)
}

// AddDefaultGroup 添加默认群组
func (s *AdminServer) AddDefaultGroup(ctx context.Context, req *admin.AddDefaultGroupReq) (*admin.AddDefaultGroupResp, error) {
	l := logic.NewAddDefaultGroupLogic(ctx, s.svcCtx)
	return l.AddDefaultGroup(req)
}

// DelDefaultGroup 删除默认群组
func (s *AdminServer) DelDefaultGroup(ctx context.Context, req *admin.DelDefaultGroupReq) (*admin.DelDefaultGroupResp, error) {
	l := logic.NewDelDefaultGroupLogic(ctx, s.svcCtx)
	return l.DelDefaultGroup(req)
}

// FindDefaultGroup 查找默认群组
func (s *AdminServer) FindDefaultGroup(ctx context.Context, req *admin.FindDefaultGroupReq) (*admin.FindDefaultGroupResp, error) {
	l := logic.NewFindDefaultGroupLogic(ctx, s.svcCtx)
	return l.FindDefaultGroup(req)
}

// SearchDefaultGroup 搜索默认群组
func (s *AdminServer) SearchDefaultGroup(ctx context.Context, req *admin.SearchDefaultGroupReq) (*admin.SearchDefaultGroupResp, error) {
	l := logic.NewSearchDefaultGroupLogic(ctx, s.svcCtx)
	return l.SearchDefaultGroup(req)
}

// AddInvitationCode 添加邀请码
func (s *AdminServer) AddInvitationCode(ctx context.Context, req *admin.AddInvitationCodeReq) (*admin.AddInvitationCodeResp, error) {
	l := logic.NewAddInvitationCodeLogic(ctx, s.svcCtx)
	return l.AddInvitationCode(req)
}

// GenInvitationCode 生成邀请码
func (s *AdminServer) GenInvitationCode(ctx context.Context, req *admin.GenInvitationCodeReq) (*admin.GenInvitationCodeResp, error) {
	l := logic.NewGenInvitationCodeLogic(ctx, s.svcCtx)
	return l.GenInvitationCode(req)
}

// FindInvitationCode 查找邀请码
func (s *AdminServer) FindInvitationCode(ctx context.Context, req *admin.FindInvitationCodeReq) (*admin.FindInvitationCodeResp, error) {
	l := logic.NewFindInvitationCodeLogic(ctx, s.svcCtx)
	return l.FindInvitationCode(req)
}

// UseInvitationCode 使用邀请码
func (s *AdminServer) UseInvitationCode(ctx context.Context, req *admin.UseInvitationCodeReq) (*admin.UseInvitationCodeResp, error) {
	l := logic.NewUseInvitationCodeLogic(ctx, s.svcCtx)
	return l.UseInvitationCode(req)
}

// DelInvitationCode 删除邀请码
func (s *AdminServer) DelInvitationCode(ctx context.Context, req *admin.DelInvitationCodeReq) (*admin.DelInvitationCodeResp, error) {
	l := logic.NewDelInvitationCodeLogic(ctx, s.svcCtx)
	return l.DelInvitationCode(req)
}

// SearchInvitationCode 搜索邀请码
func (s *AdminServer) SearchInvitationCode(ctx context.Context, req *admin.SearchInvitationCodeReq) (*admin.SearchInvitationCodeResp, error) {
	l := logic.NewSearchInvitationCodeLogic(ctx, s.svcCtx)
	return l.SearchInvitationCode(req)
}

// SearchUserIPLimitLogin 搜索用户IP登录限制
func (s *AdminServer) SearchUserIPLimitLogin(ctx context.Context, req *admin.SearchUserIPLimitLoginReq) (*admin.SearchUserIPLimitLoginResp, error) {
	l := logic.NewSearchUserIPLimitLoginLogic(ctx, s.svcCtx)
	return l.SearchUserIPLimitLogin(req)
}

// AddUserIPLimitLogin 添加用户IP登录限制
func (s *AdminServer) AddUserIPLimitLogin(ctx context.Context, req *admin.AddUserIPLimitLoginReq) (*admin.AddUserIPLimitLoginResp, error) {
	l := logic.NewAddUserIPLimitLoginLogic(ctx, s.svcCtx)
	return l.AddUserIPLimitLogin(req)
}

// DelUserIPLimitLogin 删除用户IP登录限制
func (s *AdminServer) DelUserIPLimitLogin(ctx context.Context, req *admin.DelUserIPLimitLoginReq) (*admin.DelUserIPLimitLoginResp, error) {
	l := logic.NewDelUserIPLimitLoginLogic(ctx, s.svcCtx)
	return l.DelUserIPLimitLogin(req)
}

// SearchIPForbidden 搜索IP禁止
func (s *AdminServer) SearchIPForbidden(ctx context.Context, req *admin.SearchIPForbiddenReq) (*admin.SearchIPForbiddenResp, error) {
	l := logic.NewSearchIPForbiddenLogic(ctx, s.svcCtx)
	return l.SearchIPForbidden(req)
}

// AddIPForbidden 添加IP禁止
func (s *AdminServer) AddIPForbidden(ctx context.Context, req *admin.AddIPForbiddenReq) (*admin.AddIPForbiddenResp, error) {
	l := logic.NewAddIPForbiddenLogic(ctx, s.svcCtx)
	return l.AddIPForbidden(req)
}

// DelIPForbidden 删除IP禁止
func (s *AdminServer) DelIPForbidden(ctx context.Context, req *admin.DelIPForbiddenReq) (*admin.DelIPForbiddenResp, error) {
	l := logic.NewDelIPForbiddenLogic(ctx, s.svcCtx)
	return l.DelIPForbidden(req)
}

// CancellationUser 注销用户
func (s *AdminServer) CancellationUser(ctx context.Context, req *admin.CancellationUserReq) (*admin.CancellationUserResp, error) {
	l := logic.NewCancellationUserLogic(ctx, s.svcCtx)
	return l.CancellationUser(req)
}

// BlockUser 封禁用户
func (s *AdminServer) BlockUser(ctx context.Context, req *admin.BlockUserReq) (*admin.BlockUserResp, error) {
	l := logic.NewBlockUserLogic(ctx, s.svcCtx)
	return l.BlockUser(req)
}

// UnblockUser 解封用户
func (s *AdminServer) UnblockUser(ctx context.Context, req *admin.UnblockUserReq) (*admin.UnblockUserResp, error) {
	l := logic.NewUnblockUserLogic(ctx, s.svcCtx)
	return l.UnblockUser(req)
}

// SearchBlockUser 搜索封禁用户
func (s *AdminServer) SearchBlockUser(ctx context.Context, req *admin.SearchBlockUserReq) (*admin.SearchBlockUserResp, error) {
	l := logic.NewSearchBlockUserLogic(ctx, s.svcCtx)
	return l.SearchBlockUser(req)
}

// FindUserBlockInfo 查找用户封禁信息
func (s *AdminServer) FindUserBlockInfo(ctx context.Context, req *admin.FindUserBlockInfoReq) (*admin.FindUserBlockInfoResp, error) {
	l := logic.NewFindUserBlockInfoLogic(ctx, s.svcCtx)
	return l.FindUserBlockInfo(req)
}

// CheckRegisterForbidden 检查注册是否被禁止
func (s *AdminServer) CheckRegisterForbidden(ctx context.Context, req *admin.CheckRegisterForbiddenReq) (*admin.CheckRegisterForbiddenResp, error) {
	l := logic.NewCheckRegisterForbiddenLogic(ctx, s.svcCtx)
	return l.CheckRegisterForbidden(req)
}

// CheckLoginForbidden 检查登录是否被禁止
func (s *AdminServer) CheckLoginForbidden(ctx context.Context, req *admin.CheckLoginForbiddenReq) (*admin.CheckLoginForbiddenResp, error) {
	l := logic.NewCheckLoginForbiddenLogic(ctx, s.svcCtx)
	return l.CheckLoginForbidden(req)
}

// CreateToken 创建Token
func (s *AdminServer) CreateToken(ctx context.Context, req *admin.CreateTokenReq) (*admin.CreateTokenResp, error) {
	l := logic.NewCreateTokenLogic(ctx, s.svcCtx)
	return l.CreateToken(req)
}

// ParseToken 解析Token
func (s *AdminServer) ParseToken(ctx context.Context, req *admin.ParseTokenReq) (*admin.ParseTokenResp, error) {
	l := logic.NewParseTokenLogic(ctx, s.svcCtx)
	return l.ParseToken(req)
}

// AddApplet 添加小程序
func (s *AdminServer) AddApplet(ctx context.Context, req *admin.AddAppletReq) (*admin.AddAppletResp, error) {
	l := logic.NewAddAppletLogic(ctx, s.svcCtx)
	return l.AddApplet(req)
}

// DelApplet 删除小程序
func (s *AdminServer) DelApplet(ctx context.Context, req *admin.DelAppletReq) (*admin.DelAppletResp, error) {
	l := logic.NewDelAppletLogic(ctx, s.svcCtx)
	return l.DelApplet(req)
}

// UpdateApplet 更新小程序
func (s *AdminServer) UpdateApplet(ctx context.Context, req *admin.UpdateAppletReq) (*admin.UpdateAppletResp, error) {
	l := logic.NewUpdateAppletLogic(ctx, s.svcCtx)
	return l.UpdateApplet(req)
}

// FindApplet 查找小程序
func (s *AdminServer) FindApplet(ctx context.Context, req *admin.FindAppletReq) (*admin.FindAppletResp, error) {
	l := logic.NewFindAppletLogic(ctx, s.svcCtx)
	return l.FindApplet(req)
}

// SearchApplet 搜索小程序
func (s *AdminServer) SearchApplet(ctx context.Context, req *admin.SearchAppletReq) (*admin.SearchAppletResp, error) {
	l := logic.NewSearchAppletLogic(ctx, s.svcCtx)
	return l.SearchApplet(req)
}

// GetClientConfig 获取客户端配置
func (s *AdminServer) GetClientConfig(ctx context.Context, req *admin.GetClientConfigReq) (*admin.GetClientConfigResp, error) {
	l := logic.NewGetClientConfigLogic(ctx, s.svcCtx)
	return l.GetClientConfig(req)
}

// SetClientConfig 设置客户端配置
func (s *AdminServer) SetClientConfig(ctx context.Context, req *admin.SetClientConfigReq) (*admin.SetClientConfigResp, error) {
	l := logic.NewSetClientConfigLogic(ctx, s.svcCtx)
	return l.SetClientConfig(req)
}

// DelClientConfig 删除客户端配置
func (s *AdminServer) DelClientConfig(ctx context.Context, req *admin.DelClientConfigReq) (*admin.DelClientConfigResp, error) {
	l := logic.NewDelClientConfigLogic(ctx, s.svcCtx)
	return l.DelClientConfig(req)
}

// GetUserToken 获取用户Token
func (s *AdminServer) GetUserToken(ctx context.Context, req *admin.GetUserTokenReq) (*admin.GetUserTokenResp, error) {
	l := logic.NewGetUserTokenLogic(ctx, s.svcCtx)
	return l.GetUserToken(req)
}

// InvalidateToken 使Token失效
func (s *AdminServer) InvalidateToken(ctx context.Context, req *admin.InvalidateTokenReq) (*admin.InvalidateTokenResp, error) {
	l := logic.NewInvalidateTokenLogic(ctx, s.svcCtx)
	return l.InvalidateToken(req)
}

// LatestApplicationVersion 获取最新应用版本
func (s *AdminServer) LatestApplicationVersion(ctx context.Context, req *admin.LatestApplicationVersionReq) (*admin.LatestApplicationVersionResp, error) {
	l := logic.NewLatestApplicationVersionLogic(ctx, s.svcCtx)
	return l.LatestApplicationVersion(req)
}

// AddApplicationVersion 添加应用版本
func (s *AdminServer) AddApplicationVersion(ctx context.Context, req *admin.AddApplicationVersionReq) (*admin.AddApplicationVersionResp, error) {
	l := logic.NewAddApplicationVersionLogic(ctx, s.svcCtx)
	return l.AddApplicationVersion(req)
}

// UpdateApplicationVersion 更新应用版本
func (s *AdminServer) UpdateApplicationVersion(ctx context.Context, req *admin.UpdateApplicationVersionReq) (*admin.UpdateApplicationVersionResp, error) {
	l := logic.NewUpdateApplicationVersionLogic(ctx, s.svcCtx)
	return l.UpdateApplicationVersion(req)
}

// DeleteApplicationVersion 删除应用版本
func (s *AdminServer) DeleteApplicationVersion(ctx context.Context, req *admin.DeleteApplicationVersionReq) (*admin.DeleteApplicationVersionResp, error) {
	l := logic.NewDeleteApplicationVersionLogic(ctx, s.svcCtx)
	return l.DeleteApplicationVersion(req)
}

// PageApplicationVersion 分页获取应用版本
func (s *AdminServer) PageApplicationVersion(ctx context.Context, req *admin.PageApplicationVersionReq) (*admin.PageApplicationVersionResp, error) {
	l := logic.NewPageApplicationVersionLogic(ctx, s.svcCtx)
	return l.PageApplicationVersion(req)
}
