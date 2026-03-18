package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserClientConfigLogic {
	return &DelUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserClientConfigLogic) DelUserClientConfig(req *user.DelUserClientConfigReq) (*user.DelUserClientConfigResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// TODO: 实现 DelUserConfig 方法
	return nil, fmt.Errorf("DelUserClientConfig not fully implemented yet: clientConfig database needed")
}
