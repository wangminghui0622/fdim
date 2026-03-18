package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type InitiateMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitiateMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateMultipartUploadLogic {
	return &InitiateMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateMultipartUploadLogic) InitiateMultipartUpload(req *types.InitiateMultipartUploadReq) (resp *types.InitiateMultipartUploadResp, err error) {
	// 转换请求参数
	rpcReq := &third.InitiateMultipartUploadReq{
		Hash:      req.Hash,
		Size:      req.Size,
		PartSize:  req.PartSize,
		MaxParts:  req.MaxParts,
		Cause:     req.Cause,
		UrlPrefix: req.UrlPrefix,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.InitiateMultipartUpload(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.InitiateMultipartUploadResp{
		Url:    rpcResp.Url,
		Upload: rpcResp.Upload,
	}

	return resp, nil
}
