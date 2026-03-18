package logic

import (
	"context"
	"fmt"

	"fdim/protocol/auth"
	"fdim/protocol/constant"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTokenLogic {
	return &GetUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserTokenLogic) GetUserToken(req *auth.GetUserTokenReq) (*auth.GetUserTokenResp, error) {
	// 注：移除管理员权限检查，允许用户自行获取token
	// if err := authverify.CheckAdmin(l.ctx); err != nil {
	// 	return nil, err
	// }

	// 检查 PlatformID
	if req.PlatformID == constant.AdminPlatformID {
		return nil, fmt.Errorf("platformID invalid. platformID must not be adminPlatformID")
	}

	if l.svcCtx.AuthDB == nil {
		return nil, fmt.Errorf("auth database not initialized")
	}

	// 验证 userID 对应的用户是否存在（防止用昵称等无效ID获取token）
	if l.svcCtx.UserDB != nil {
		users, err := l.svcCtx.UserDB.FindWithError(l.ctx, []string{req.UserID})
		if err != nil {
			l.Errorf("GetUserToken: query user failed: userID=%s, err=%v", req.UserID, err)
			return nil, fmt.Errorf("user %s not found", req.UserID)
		}
		if len(users) == 0 {
			l.Errorf("GetUserToken: user not found: userID=%s", req.UserID)
			return nil, fmt.Errorf("user %s not found", req.UserID)
		}
	}

	// 2. 创建用户 Token（UserType=1 表示普通用户）
	tokenStr, expireSeconds, err := l.svcCtx.AuthDB.CreateToken(l.ctx, req.UserID, req.PlatformID, 1)
	if err != nil {
		return nil, err
	}

	return &auth.GetUserTokenResp{
		Token:             tokenStr,
		ExpireTimeSeconds: expireSeconds,
	}, nil
}
