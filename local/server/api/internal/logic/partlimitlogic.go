package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type PartLimitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPartLimitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PartLimitLogic {
	return &PartLimitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PartLimitLogic) PartLimit(req *types.PartLimitReq) (resp *types.PartLimitResp, err error) {
	// 转换请求参数
	rpcReq := &third.PartLimitReq{}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.PartLimit(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.PartLimitResp{
		MinPartSize: rpcResp.MinPartSize,
		MaxPartSize: rpcResp.MaxPartSize,
		MaxNumSize:  rpcResp.MaxNumSize,
	}

	return resp, nil
}
