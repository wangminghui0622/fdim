package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type CompleteMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(req *types.CompleteMultipartUploadReq) (resp *types.CompleteMultipartUploadResp, err error) {
	// 转换请求参数
	rpcReq := &third.CompleteMultipartUploadReq{
		UploadID:    req.UploadID,
		Parts:       req.Parts,
		Name:        req.Name,
		ContentType: req.ContentType,
		Cause:       req.Cause,
		UrlPrefix:   req.UrlPrefix,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.CompleteMultipartUpload(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.CompleteMultipartUploadResp{
		Url: rpcResp.Url,
	}

	return resp, nil
}
