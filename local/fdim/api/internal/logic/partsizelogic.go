package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type PartSizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPartSizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PartSizeLogic {
	return &PartSizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PartSizeLogic) PartSize(req *types.PartSizeReq) (resp *types.PartSizeResp, err error) {
	// 转换请求参数
	rpcReq := &third.PartSizeReq{
		Size: req.Size,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.PartSize(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.PartSizeResp{
		Size: rpcResp.Size,
	}

	return resp, nil
}
