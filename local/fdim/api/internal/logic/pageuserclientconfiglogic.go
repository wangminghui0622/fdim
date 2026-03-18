package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type PageUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageUserClientConfigLogic {
	return &PageUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageUserClientConfigLogic) PageUserClientConfig(req *types.PageUserClientConfigReq) (resp *types.PageUserClientConfigResp, err error) {
	// 转换请求参数
	rpcReq := &user.PageUserClientConfigReq{
		UserID: req.UserID,
		Key:    req.Key,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.PageUserClientConfig(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var configs []types.ClientConfig
	for _, cfg := range rpcResp.Configs {
		configs = append(configs, types.ClientConfig{
			Key:   cfg.Key,
			Value: cfg.Value,
		})
	}

	resp = &types.PageUserClientConfigResp{
		Total:   rpcResp.Total,
		Configs: configs,
	}

	return resp, nil
}
