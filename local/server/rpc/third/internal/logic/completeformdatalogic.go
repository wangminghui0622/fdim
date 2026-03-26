package logic

import (
	"context"
	"fmt"
	"strings"

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

	// 从 uploadID 中解码对象名（格式: "form:objectName"）
	objectName := ""
	if parts := strings.SplitN(req.Id, ":", 2); len(parts) == 2 {
		objectName = parts[1]
	}
	if objectName == "" {
		return nil, fmt.Errorf("invalid upload id: %s", req.Id)
	}

	// 用 SignEndpoint（公网地址）构造访问 URL，回退到 Endpoint
	endpoint := cfg.SignEndpoint
	if endpoint == "" {
		endpoint = cfg.Endpoint
	}
	endpoint = strings.TrimRight(endpoint, "/")

	accessUrl := fmt.Sprintf("%s/%s/%s", endpoint, cfg.Bucket, objectName)

	return &third.CompleteFormDataResp{
		Url: accessUrl,
	}, nil
}
