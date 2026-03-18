package logic

import (
	"context"
	"fmt"
	"net/url"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthSignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthSignLogic {
	return &AuthSignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthSignLogic) AuthSign(req *third.AuthSignReq) (*third.AuthSignResp, error) {
	// 简化版：不生成真正的 S3 签名，而是基于 UploadID 和 PartNumbers 拼接上传 URL 和部分 query/header
	if req.UploadID == "" {
		return nil, fmt.Errorf("uploadID is empty")
	}
	if len(req.PartNumbers) == 0 {
		return nil, fmt.Errorf("partNumbers is empty")
	}

	cfg := l.svcCtx.Config.ObjectStorage
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("object storage not configured")
	}

	// 这里使用 UploadID 作为对象路径前缀，客户端按需使用
	baseURL := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, cfg.Bucket, url.PathEscape(req.UploadID))
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upload url: %w", err)
	}

	// 将 partNumbers 以逗号拼接到 query 中，方便客户端识别
	q := u.Query()
	for _, part := range req.PartNumbers {
		q.Add("partNumber", fmt.Sprintf("%d", part))
	}
	u.RawQuery = q.Encode()

	return &third.AuthSignResp{
		Url:    u.String(),
		Query:  []*third.KeyValues{}, // 暂不生成额外签名参数
		Header: []*third.KeyValues{},
		Parts:  []*third.SignPart{},
	}, nil
}
