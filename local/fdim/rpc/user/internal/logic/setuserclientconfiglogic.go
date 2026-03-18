package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserClientConfigLogic {
	return &SetUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserClientConfigLogic) SetUserClientConfig(req *user.SetUserClientConfigReq) (*user.SetUserClientConfigResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	if req.UserID != "" {
		if _, err := l.svcCtx.UserDB.GetUserByID(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// TODO: 实现 SetUserConfig 方法
	return nil, fmt.Errorf("SetUserClientConfig not fully implemented yet: clientConfig database needed")
}
