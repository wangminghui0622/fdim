package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type AuthSignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthSignLogic {
	return &AuthSignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthSignLogic) AuthSign(req *types.AuthSignReq) (resp *types.AuthSignResp, err error) {
	// 转换请求参数
	rpcReq := &third.AuthSignReq{
		UploadID:    req.UploadID,
		PartNumbers: req.PartNumbers,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.AuthSign(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	query := make(map[string]string)
	for _, kv := range rpcResp.Query {
		if len(kv.Values) > 0 {
			query[kv.Key] = kv.Values[0]
		}
	}

	header := make(map[string]string)
	for _, kv := range rpcResp.Header {
		if len(kv.Values) > 0 {
			header[kv.Key] = kv.Values[0]
		}
	}

	resp = &types.AuthSignResp{
		Url:    rpcResp.Url,
		Query:  query,
		Header: header,
	}

	return resp, nil
}
