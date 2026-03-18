package logic

import (
	"context"
	"fmt"
	"net/url"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	// 简化版：不真正向对象存储服务发送 Complete 请求，而是基于配置拼出最终访问 URL
	if err := req.Check(); err != nil {
		return nil, fmt.Errorf("invalid CompleteMultipartUploadReq: %w", err)
	}

	cfg := l.svcCtx.Config.ObjectStorage
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("object storage not configured")
	}

	name := req.Name
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	baseURL := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, cfg.Bucket, url.PathEscape(name))
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object url: %w", err)
	}

	return &third.CompleteMultipartUploadResp{
		Url: u.String(),
	}, nil
}
