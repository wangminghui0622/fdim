package server

import (
	"context"

	"fdim/protocol/auth"
	"fdim/rpc/auth/internal/logic"
	"fdim/rpc/auth/internal/svc"
)

type AuthServer struct {
	auth.UnimplementedAuthServer
	svcCtx *svc.ServiceContext
}

func NewAuthServer(svcCtx *svc.ServiceContext) *AuthServer {
	return &AuthServer{
		svcCtx: svcCtx,
	}
}

// GetAdminToken 获取管理员Token
func (s *AuthServer) GetAdminToken(ctx context.Context, req *auth.GetAdminTokenReq) (*auth.GetAdminTokenResp, error) {
	l := logic.NewGetAdminTokenLogic(ctx, s.svcCtx)
	return l.GetAdminToken(req)
}

// GetUserToken 获取用户Token
func (s *AuthServer) GetUserToken(ctx context.Context, req *auth.GetUserTokenReq) (*auth.GetUserTokenResp, error) {
	l := logic.NewGetUserTokenLogic(ctx, s.svcCtx)
	return l.GetUserToken(req)
}

// ForceLogout 强制登出
func (s *AuthServer) ForceLogout(ctx context.Context, req *auth.ForceLogoutReq) (*auth.ForceLogoutResp, error) {
	l := logic.NewForceLogoutLogic(ctx, s.svcCtx)
	return l.ForceLogout(req)
}

// ParseToken 解析Token
func (s *AuthServer) ParseToken(ctx context.Context, req *auth.ParseTokenReq) (*auth.ParseTokenResp, error) {
	l := logic.NewParseTokenLogic(ctx, s.svcCtx)
	return l.ParseToken(req)
}

// InvalidateToken 使Token失效
func (s *AuthServer) InvalidateToken(ctx context.Context, req *auth.InvalidateTokenReq) (*auth.InvalidateTokenResp, error) {
	l := logic.NewInvalidateTokenLogic(ctx, s.svcCtx)
	return l.InvalidateToken(req)
}

// KickTokens 踢出Token
func (s *AuthServer) KickTokens(ctx context.Context, req *auth.KickTokensReq) (*auth.KickTokensResp, error) {
	l := logic.NewKickTokensLogic(ctx, s.svcCtx)
	return l.KickTokens(req)
}

// GetExistingToken 获取现有Token
func (s *AuthServer) GetExistingToken(ctx context.Context, req *auth.GetExistingTokenReq) (*auth.GetExistingTokenResp, error) {
	l := logic.NewGetExistingTokenLogic(ctx, s.svcCtx)
	return l.GetExistingToken(req)
}
