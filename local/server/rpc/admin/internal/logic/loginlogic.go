package logic

import (
	"context"
	"fmt"

	"fdim/pkg/constant"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *admin.LoginReq) (*admin.LoginResp, error) {
	// 1. 从数据库获取管理员信息
	adminInfo, err := l.svcCtx.AdminDB.GetAdmin(l.ctx, req.Account)
	if err != nil {
		l.Errorf("GetAdmin failed: %v", err)
		return nil, fmt.Errorf("account not found")
	}

	// 2. 验证密码
	if adminInfo.Password != req.Password {
		return nil, fmt.Errorf("invalid password")
	}

	// 3. 生成 Token
	tokenStr, expireDuration, err := l.svcCtx.Token.CreateToken(adminInfo.UserID, constant.AdminUser)
	if err != nil {
		l.Errorf("CreateToken failed: %v", err)
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	// 4. 缓存 Token 到 Redis
	expireSeconds := int64(expireDuration.Seconds())
	if err := l.svcCtx.AdminDB.CacheToken(l.ctx, adminInfo.UserID, tokenStr, expireSeconds); err != nil {
		l.Errorf("CacheToken failed: %v", err)
		// Token 已生成，即使缓存失败也返回成功
	}

	// 5. 返回登录结果
	return &admin.LoginResp{
		AdminUserID:  adminInfo.UserID,
		AdminAccount: adminInfo.Account,
		AdminToken:   tokenStr,
		Nickname:     adminInfo.Nickname,
		FaceURL:      adminInfo.FaceURL,
		Level:        adminInfo.Level,
	}, nil
}
