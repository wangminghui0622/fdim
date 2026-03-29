package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClientConfigLogic {
	return &GetClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetClientConfigLogic) GetClientConfig(req *admin.GetClientConfigReq) (*admin.GetClientConfigResp, error) {
	config, err := l.svcCtx.AdminDB.GetConfig(l.ctx)
	if err != nil {
		l.Errorf("GetConfig failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to get client config")
	}

	return &admin.GetClientConfigResp{
		Config: config,
	}, nil
}
