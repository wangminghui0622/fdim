package logic

import (
	"context"
	"fmt"
	"net/url"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteFormDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteFormDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteFormDataLogic {
	return &CompleteFormDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteFormDataLogic) CompleteFormData(req *third.CompleteFormDataReq) (*third.CompleteFormDataResp, error) {
	if err := req.Check(); err != nil {
		return nil, fmt.Errorf("invalid CompleteFormDataReq: %w", err)
	}

	cfg := l.svcCtx.Config.ObjectStorage
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("object storage not configured")
	}

	// 简化版：完成表单上传后，返回一个前缀 URL，用于前端拼接最终对象路径
	baseURL := fmt.Sprintf("%s/%s", cfg.Endpoint, cfg.Bucket)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form complete url: %w", err)
	}

	return &third.CompleteFormDataResp{
		Url: u.String(),
	}, nil
}
