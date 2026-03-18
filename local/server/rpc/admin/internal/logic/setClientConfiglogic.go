package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetClientConfigLogic {
	return &SetClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetClientConfigLogic) SetClientConfig(req *admin.SetClientConfigReq) (*admin.SetClientConfigResp, error) {
	// 1. 验证参数
	if len(req.Config) == 0 {
		return nil, errs.ErrArgs.WrapMsg("config cannot be empty")
	}

	// 2. 设置配置
	if err := l.svcCtx.AdminDB.SetConfig(l.ctx, req.Config); err != nil {
		l.Errorf("SetConfig failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to set client config")
	}

	l.Infof("Set client config: keys=%d", len(req.Config))
	return &admin.SetClientConfigResp{}, nil
}
