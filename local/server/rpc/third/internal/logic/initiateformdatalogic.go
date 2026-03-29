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

type InitiateFormDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitiateFormDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateFormDataLogic {
	return &InitiateFormDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitiateFormDataLogic) InitiateFormData(req *third.InitiateFormDataReq) (*third.InitiateFormDataResp, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	// 检查对象存储是否已配置
	if l.svcCtx.ObjectStorage != nil {
		// 使用真实?MinIO FormData 接口
		formData, err := l.svcCtx.ObjectStorage.FormData(l.ctx, req.Name, req.Size, req.ContentType, 10*time.Minute)
		if err != nil {
			l.Logger.Errorf("FormData failed: %v", err)
			return nil, fmt.Errorf("failed to get form data: %w", err)
		}

		// 把对象名编码?uploadID，CompleteFormData 用它来构造访?URL
		uploadID := fmt.Sprintf("form:%s", req.Name)
		successCodes := make([]int32, len(formData.SuccessCodes))
		for i, code := range formData.SuccessCodes {
			successCodes[i] = int32(code)
		}

		return &third.InitiateFormDataResp{
			Id:           uploadID,
			Url:          formData.URL,
			File:         formData.File,
			Header:       []*third.KeyValues{},
			FormData:     formData.FormData,
			Expires:      formData.Expires.Unix(),
			SuccessCodes: successCodes,
		}, nil
	}

	// 简化版逻辑
	cfg := l.svcCtx.Config.ObjectStorage
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("object storage not configured")
	}

	uploadID := fmt.Sprintf("form-%d", time.Now().UnixNano())
	baseURL := fmt.Sprintf("%s/%s", cfg.Endpoint, cfg.Bucket)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form upload url: %w", err)
	}
	q := u.Query()
	q.Set("uploadID", uploadID)
	u.RawQuery = q.Encode()

	return &third.InitiateFormDataResp{
		Id:           uploadID,
		Url:          u.String(),
		File:         req.Name,
		Header:       []*third.KeyValues{},
		FormData:     map[string]string{},
		Expires:      time.Now().Add(10 * time.Minute).Unix(),
		SuccessCodes: []int32{200, 204},
	}, nil
}
