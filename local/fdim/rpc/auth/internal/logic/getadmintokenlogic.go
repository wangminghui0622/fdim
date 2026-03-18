package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/auth"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminTokenLogic {
	return &GetAdminTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminTokenLogic) GetAdminToken(req *auth.GetAdminTokenReq) (*auth.GetAdminTokenResp, error) {
	// 1. 验证请求中的 secret 是否与配置一致
	if req.Secret == "" || req.Secret != l.svcCtx.Config.Secret {
		return nil, fmt.Errorf("invalid secret")
	}
	// 2. 验证 userID 是否为管理员
	if !authverify.CheckUserIsAdmin(l.ctx, req.UserID) {
		return nil, fmt.Errorf("user %s is not admin", req.UserID)
	}
	// 3. 创建管理员 Token（UserType=Admin）
	if l.svcCtx.Token == nil || l.svcCtx.AuthDB == nil {
		return nil, fmt.Errorf("auth service not initialized")
	}
	tokenStr, expireSeconds, err := l.svcCtx.AuthDB.CreateToken(l.ctx, req.UserID, 0, // 管理端不区分平台ID
		1) // 使用 admin 类型，在 Token 内部做区分
	if err != nil {
		return nil, err
	}

	return &auth.GetAdminTokenResp{
		Token:             tokenStr,
		ExpireTimeSeconds: expireSeconds,
	}, nil
}
