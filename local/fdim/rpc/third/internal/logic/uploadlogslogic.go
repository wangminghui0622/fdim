package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
)

type UploadLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogsLogic {
	return &UploadLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadLogsLogic) UploadLogs(req *third.UploadLogsReq) (*third.UploadLogsResp, error) {
	if l.svcCtx.MongoDB == nil {
		return nil, fmt.Errorf("mongo not initialized")
	}
	if len(req.FileURLs) == 0 {
		return &third.UploadLogsResp{}, nil
	}

	coll := l.svcCtx.MongoDB.GetCollection("third_logs")
	now := time.Now().UnixMilli()
	docs := make([]interface{}, 0, len(req.FileURLs))
	for _, f := range req.FileURLs {
		docs = append(docs, bson.M{
			"platform":     req.Platform,
			"url":          f.URL,
			"filename":     f.Filename,
			"appFramework": req.AppFramework,
			"version":      req.Version,
			"ex":           req.Ex,
			"createTime":   now,
		})
	}

	if _, err := coll.InsertMany(l.ctx, docs); err != nil {
		l.Errorf("failed to insert logs: %v", err)
		return nil, fmt.Errorf("failed to save logs: %w", err)
	}

	return &third.UploadLogsResp{}, nil
}
