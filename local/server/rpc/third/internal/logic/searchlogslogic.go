package logic

import (
	"context"
	"fmt"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogsLogic {
	return &SearchLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchLogsLogic) SearchLogs(req *third.SearchLogsReq) (*third.SearchLogsResp, error) {
	if l.svcCtx.MongoDB == nil {
		return nil, fmt.Errorf("MongoDB未初始化")
	}

	coll := l.svcCtx.MongoDB.GetCollection("third_logs")

	filter := bson.M{}
	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"filename": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"url": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"ex": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}
	if req.StartTime > 0 || req.EndTime > 0 {
		timeCond := bson.M{}
		if req.StartTime > 0 {
			timeCond["$gte"] = req.StartTime
		}
		if req.EndTime > 0 {
			timeCond["$lte"] = req.EndTime
		}
		filter["createTime"] = timeCond
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	if req.Pagination != nil {
		skip := int64((req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber)
		limit := int64(req.Pagination.ShowNumber)
		if skip < 0 {
			skip = 0
		}
		if limit <= 0 {
			limit = 20
		}
		opts.SetSkip(skip).SetLimit(limit)
	}

	cursor, err := coll.Find(l.ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询日志失败: %w", err)
	}
	defer cursor.Close(l.ctx)

	logsInfos := make([]*third.LogInfo, 0)
	for cursor.Next(l.ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("解码日志失败: %w", err)
		}
		li := &third.LogInfo{}
		if v, ok := doc["userID"].(string); ok {
			li.UserID = v
		}
		if v, ok := doc["platform"].(int32); ok {
			li.Platform = fmt.Sprintf("%d", v)
		} else if v, ok := doc["platform"].(int64); ok {
			li.Platform = fmt.Sprintf("%d", v)
		}
		if v, ok := doc["url"].(string); ok {
			li.Url = v
		}
		if v, ok := doc["createTime"].(int64); ok {
			li.CreateTime = v
		}
		if v, ok := doc["nickname"].(string); ok {
			li.Nickname = v
		}
		if v, ok := doc["logID"].(string); ok {
			li.LogID = v
		}
		if v, ok := doc["filename"].(string); ok {
			li.Filename = v
		}
		if v, ok := doc["systemType"].(string); ok {
			li.SystemType = v
		}
		if v, ok := doc["ex"].(string); ok {
			li.Ex = v
		}
		if v, ok := doc["version"].(string); ok {
			li.Version = v
		}
		logsInfos = append(logsInfos, li)
	}

	total, err := coll.CountDocuments(l.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计日志数量失败: %w", err)
	}

	return &third.SearchLogsResp{
		LogsInfos: logsInfos,
		Total:     uint32(total),
	}, nil
}
