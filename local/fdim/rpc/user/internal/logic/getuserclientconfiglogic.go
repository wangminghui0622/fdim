package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserClientConfigLogic {
	return &GetUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserClientConfigLogic) GetUserClientConfig(req *user.GetUserClientConfigReq) (*user.GetUserClientConfigResp, error) {
	// 权限验证
	if req.UserID != "" {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
		if _, err := l.svcCtx.UserDB.GetUserByID(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// TODO: 实现 GetUserConfig 方法
	return nil, fmt.Errorf("GetUserClientConfig not fully implemented yet: clientConfig database needed")
}
