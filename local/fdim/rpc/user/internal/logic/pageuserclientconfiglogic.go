package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PageUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageUserClientConfigLogic {
	return &PageUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageUserClientConfigLogic) PageUserClientConfig(req *user.PageUserClientConfigReq) (*user.PageUserClientConfigResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// TODO: 实现 GetUserConfigPage 方法
	return nil, fmt.Errorf("PageUserClientConfig not fully implemented yet: clientConfig database needed")
}
