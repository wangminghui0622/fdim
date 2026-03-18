package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type InitiateFormDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitiateFormDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateFormDataLogic {
	return &InitiateFormDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateFormDataLogic) InitiateFormData(req *types.InitiateFormDataReq) (resp *types.InitiateFormDataResp, err error) {
	// 转换请求参数
	rpcReq := &third.InitiateFormDataReq{
		Name:        req.Name,
		Size:        req.Size,
		ContentType: req.ContentType,
		Group:       req.Group,
		Millisecond: req.Millisecond,
		Absolute:    req.Absolute,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.InitiateFormData(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var header []interface{}
	for _, kv := range rpcResp.Header {
		header = append(header, kv)
	}

	resp = &types.InitiateFormDataResp{
		Id:           rpcResp.Id,
		Url:          rpcResp.Url,
		File:         rpcResp.File,
		Header:       header,
		FormData:     rpcResp.FormData,
		Expires:      rpcResp.Expires,
		SuccessCodes: rpcResp.SuccessCodes,
	}

	return resp, nil
}
