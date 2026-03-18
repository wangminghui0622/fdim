package logic

import (
	"context"
	"time"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetActiveGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActiveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveGroupLogic {
	return &GetActiveGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetActiveGroupLogic) GetActiveGroup(req *msg.GetActiveGroupReq) (*msg.GetActiveGroupResp, error) {
	if req.End == 0 {
		req.End = time.Now().UnixMilli()
	}
	if req.Start == 0 || req.Start > req.End {
		// 默认统计最近 7 天
		req.Start = req.End - int64(7*24*time.Hour/time.Millisecond)
	}
	if req.Pagination == nil || req.Pagination.ShowNumber <= 0 {
		req.Pagination = &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 20,
		}
	}

	coll := l.svcCtx.MongoDB.GetCollection("stream_msg")
	if coll == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message collection not initialized")
	}

	match := bson.M{
		"send_time": bson.M{"$gte": req.Start, "$lte": req.End},
		"group_id":  bson.M{"$ne": ""},
	}

	// 1. 统计总消息数和活跃群数量
	totalMsgCount, err := coll.CountDocuments(l.ctx, match)
	if err != nil {
		l.Errorw("CountDocuments total msg failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to count group messages")
	}

	groupCountAgg := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{"_id": "$group_id"}}},
		{{Key: "$count", Value: "groupCount"}},
	}
	cur, err := coll.Aggregate(l.ctx, groupCountAgg)
	if err != nil {
		l.Errorw("Aggregate groupCount failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to count active groups")
	}
	var groupCount int64
	if cur.Next(l.ctx) {
		var res struct {
			GroupCount int64 `bson:"groupCount"`
		}
		if err := cur.Decode(&res); err == nil {
			groupCount = res.GroupCount
		}
	}
	cur.Close(l.ctx)

	// 2. 按天统计 dateCount
	dateAgg := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$project", Value: bson.M{
			"day": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   bson.M{"$toDate": "$send_time"},
				},
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$day",
			"count": bson.M{"$sum": 1},
		}}},
	}
	cur, err = coll.Aggregate(l.ctx, dateAgg)
	if err != nil {
		l.Errorw("Aggregate group dateCount failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to aggregate group date count")
	}
	dateCount := make(map[string]int64)
	for cur.Next(l.ctx) {
		var res struct {
			Day   string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cur.Decode(&res); err == nil {
			dateCount[res.Day] = res.Count
		}
	}
	cur.Close(l.ctx)

	// 3. 按 group 聚合消息数量，做分页
	offset := int64((req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber)
	limit := int64(req.Pagination.ShowNumber)
	groupAgg := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$group_id",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
	}
	cur, err = coll.Aggregate(l.ctx, groupAgg)
	if err != nil {
		l.Errorw("Aggregate active groups failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to aggregate active groups")
	}
	type groupAggRes struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	var groupsAgg []groupAggRes
	if err := cur.All(l.ctx, &groupsAgg); err != nil {
		l.Errorw("Decode active groups failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to decode active groups")
	}

	activeGroups := make([]*msg.ActiveGroup, 0, len(groupsAgg))
	for _, g := range groupsAgg {
		activeGroups = append(activeGroups, &msg.ActiveGroup{
			Group: &sdkws.GroupInfo{
				GroupID: g.ID,
			},
			Count: g.Count,
		})
	}

	return &msg.GetActiveGroupResp{
		MsgCount:   totalMsgCount,
		GroupCount: groupCount,
		DateCount:  dateCount,
		Groups:     activeGroups,
	}, nil
}
