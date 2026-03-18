package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type ObjectRedirectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewObjectRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ObjectRedirectLogic {
	return &ObjectRedirectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ObjectRedirectLogic) ObjectRedirect(req *types.AccessURLReq) (resp *types.AccessURLResp, err error) {
	// 转换请求参数
	rpcReq := &third.AccessURLReq{
		Name:  req.Name,
		Query: req.Query,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.AccessURL(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.AccessURLResp{
		Url: rpcResp.Url,
	}

	return resp, nil
}
