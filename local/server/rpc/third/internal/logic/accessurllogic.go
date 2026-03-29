package logic

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AccessURLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessURLLogic {
	return &AccessURLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessURLLogic) AccessURL(req *third.AccessURLReq) (*third.AccessURLResp, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	// 使用真实?MinIO 接口
	if l.svcCtx.ObjectStorage != nil {
		accessURL, err := l.svcCtx.ObjectStorage.AccessURL(l.ctx, req.Name, time.Hour, nil)
		if err != nil {
			l.Logger.Errorf("AccessURL failed: %v", err)
			return nil, fmt.Errorf("failed to get access url: %w", err)
		}

		return &third.AccessURLResp{
			Url:        accessURL,
			ExpireTime: time.Now().Add(time.Hour).Unix(),
		}, nil
	}

	// 简化版逻辑
	cfg := l.svcCtx.Config.ObjectStorage
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("object storage not configured")
	}

	baseURL := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, cfg.Bucket, url.PathEscape(req.Name))
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}
	if len(req.Query) > 0 {
		q := u.Query()
		for k, v := range req.Query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	expireTime := time.Now().Add(time.Hour).Unix()
	q := u.Query()
	q.Set("expire", strconv.FormatInt(expireTime, 10))
	u.RawQuery = q.Encode()

	return &third.AccessURLResp{
		Url:        u.String(),
		ExpireTime: expireTime,
	}, nil
}
