package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelClientConfigLogic {
	return &DelClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelClientConfigLogic) DelClientConfig(req *admin.DelClientConfigReq) (*admin.DelClientConfigResp, error) {
	// 1. 验证参数
	if len(req.Keys) == 0 {
		return nil, errs.ErrArgs.WrapMsg("keys cannot be empty")
	}

	// 2. 删除配置
	if err := l.svcCtx.AdminDB.DelConfig(l.ctx, req.Keys); err != nil {
		l.Errorf("DelConfig failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete client config")
	}

	l.Infof("Deleted client config: keys=%v", req.Keys)
	return &admin.DelClientConfigResp{}, nil
}
