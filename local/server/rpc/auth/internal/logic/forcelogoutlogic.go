package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/auth"
	"fdim/protocol/msggateway"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ForceLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewForceLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutLogic {
	return &ForceLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ForceLogoutLogic) ForceLogout(req *auth.ForceLogoutReq) (*auth.ForceLogoutResp, error) {
	l.Infof("[Force Logout] Request received - userID=%s, platformID=%d", req.UserID, req.PlatformID)
	
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	if l.svcCtx.AuthDB == nil {
		return nil, fmt.Errorf("auth database not initialized")
	}

	// 1. 在 AuthDB 中将该用户在指定平台上的所有 token 标记为 Kicked
	statusMap, err := l.svcCtx.AuthDB.GetTokens(l.ctx, req.UserID, req.PlatformID)
	if err != nil {
		return nil, err
	}
	var tokens []string
	for tk := range statusMap {
		tokens = append(tokens, tk)
	}
	if err := l.svcCtx.AuthDB.BatchKickTokens(l.ctx, req.UserID, req.PlatformID, tokens); err != nil {
		return nil, err
	}
	
	l.Infof("[Force Logout] Tokens marked as kicked - userID=%s, platformID=%d, tokenCount=%d",
		req.UserID, req.PlatformID, len(tokens))

	// 2. 调用 MsgGateway 踢出用户 WebSocket 连接（如果已配置）
	if l.svcCtx.MsgGatewayClient != nil {
		l.Infof("[Force Logout] Calling MsgGateway to kick user - userID=%s, platformID=%d",
			req.UserID, req.PlatformID)
		_, err = l.svcCtx.MsgGatewayClient.KickUserOffline(l.ctx, &msggateway.KickUserOfflineReq{
			PlatformID:     req.PlatformID,
			KickUserIDList: []string{req.UserID},
		})
		if err != nil {
			l.Errorf("[Force Logout] Failed to kick user from MsgGateway - userID=%s, platformID=%d, error=%v",
				req.UserID, req.PlatformID, err)
		}
	}

	l.Infof("[Force Logout] Force logout completed successfully - userID=%s, platformID=%d",
		req.UserID, req.PlatformID)
	return &auth.ForceLogoutResp{}, nil
}
