package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
)

type GetPrometheusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPrometheusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrometheusLogic {
	return &GetPrometheusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPrometheusLogic) GetPrometheus() (resp *types.GetPrometheusResp, err error) {
	// GetPrometheus 直接返回配置中的 Grafana URL
	grafanaURL := "" // 从配置中获取，暂时为?
	if grafanaURL == "" {
		grafanaURL = "http://localhost:3000" // 默认?
	}

	resp = &types.GetPrometheusResp{
		Url: grafanaURL,
	}

	return resp, nil
}
