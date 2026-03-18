package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type CompleteFormDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteFormDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteFormDataLogic {
	return &CompleteFormDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteFormDataLogic) CompleteFormData(req *types.CompleteFormDataReq) (resp *types.CompleteFormDataResp, err error) {
	// 转换请求参数
	rpcReq := &third.CompleteFormDataReq{
		Id:        req.Id,
		UrlPrefix: req.UrlPrefix,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.CompleteFormData(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.CompleteFormDataResp{
		Url: rpcResp.Url,
	}

	return resp, nil
}
