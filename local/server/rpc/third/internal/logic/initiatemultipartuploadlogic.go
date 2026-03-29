package logic

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InitiateMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitiateMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateMultipartUploadLogic {
	return &InitiateMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitiateMultipartUploadLogic) InitiateMultipartUpload(req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	// 参数校验
	if req.Hash == "" && req.Name == "" {
		return nil, fmt.Errorf("name and hash are both empty")
	}

	// 检查对象存储是否已配置
	if l.svcCtx.ObjectStorage == nil {
		cfg := l.svcCtx.Config.ObjectStorage
		if cfg.Endpoint == "" || cfg.Bucket == "" {
			return nil, fmt.Errorf("object storage not configured")
		}
		// 使用简化版逻辑
		return l.initiateMultipartUploadSimple(req)
	}

	// 使用真实?MinIO 接口
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("FDIM/data/%s/%s", req.Hash[:2], req.Hash)
	}

	result, err := l.svcCtx.ObjectStorage.InitiateMultipartUpload(l.ctx, name, nil)
	if err != nil {
		l.Logger.Errorf("InitiateMultipartUpload failed: %v", err)
		return nil, fmt.Errorf("failed to initiate multipart upload: %w", err)
	}

	partSize := req.PartSize
	if partSize <= 0 {
		partSize = 5 * 1024 * 1024
	}

	uploadInfo := &third.UploadInfo{
		UploadID:   result.UploadID,
		PartSize:   partSize,
		Sign:       &third.AuthSignParts{Url: "", Query: []*third.KeyValues{}, Header: []*third.KeyValues{}, Parts: []*third.SignPart{}},
		ExpireTime: time.Now().Add(10 * time.Minute).Unix(),
	}

	return &third.InitiateMultipartUploadResp{
		Url:    fmt.Sprintf("%s/%s/%s", l.svcCtx.Config.ObjectStorage.Endpoint, result.Bucket, result.Key),
		Upload: uploadInfo,
	}, nil
}

func (l *InitiateMultipartUploadLogic) initiateMultipartUploadSimple(req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	cfg := l.svcCtx.Config.ObjectStorage

	name := req.Name
	if name == "" {
		name = req.Hash
	}

	uploadID := fmt.Sprintf("%s-%d", req.Hash, time.Now().UnixNano())
	baseURL := fmt.Sprintf("%s/%s/%s", cfg.Endpoint, cfg.Bucket, url.PathEscape(name))
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse upload base url: %w", err)
	}
	q := u.Query()
	q.Set("uploadID", uploadID)
	u.RawQuery = q.Encode()

	partSize := req.PartSize
	if partSize <= 0 {
		partSize = 5 * 1024 * 1024
	}
	uploadInfo := &third.UploadInfo{
		UploadID:   uploadID,
		PartSize:   partSize,
		Sign:       &third.AuthSignParts{Url: u.String(), Query: []*third.KeyValues{}, Header: []*third.KeyValues{}, Parts: []*third.SignPart{}},
		ExpireTime: time.Now().Add(10 * time.Minute).Unix(),
	}

	return &third.InitiateMultipartUploadResp{
		Url:    u.String(),
		Upload: uploadInfo,
	}, nil
}
