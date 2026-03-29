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

type GetActiveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActiveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveUserLogic {
	return &GetActiveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetActiveUserLogic) GetActiveUser(req *msg.GetActiveUserReq) (*msg.GetActiveUserResp, error) {
	if req.End == 0 {
		req.End = time.Now().UnixMilli()
	}
	if req.Start == 0 || req.Start > req.End {
		// 默认统计最?7 ?
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
	}
	if req.Group {
		// 仅统计群聊消?
		match["group_id"] = bson.M{"$ne": ""}
	}

	// 1. 统计总消息数和活跃用户数
	totalMsgCount, err := coll.CountDocuments(l.ctx, match)
	if err != nil {
		l.Errorw("CountDocuments total msg failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to count messages")
	}

	userCountAgg := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{"_id": "$send_id"}}},
		{{Key: "$count", Value: "userCount"}},
	}
	cur, err := coll.Aggregate(l.ctx, userCountAgg)
	if err != nil {
		l.Errorw("Aggregate userCount failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to count active users")
	}
	var userCount int64
	if cur.Next(l.ctx) {
		var res struct {
			UserCount int64 `bson:"userCount"`
		}
		if err := cur.Decode(&res); err == nil {
			userCount = res.UserCount
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
		l.Errorw("Aggregate dateCount failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to aggregate date count")
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

	// 3. ?user 聚合消息数量，做分页
	offset := int64((req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber)
	limit := int64(req.Pagination.ShowNumber)
	userAgg := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$send_id",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
	}
	cur, err = coll.Aggregate(l.ctx, userAgg)
	if err != nil {
		l.Errorw("Aggregate active users failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to aggregate active users")
	}
	type userAggRes struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	var usersAgg []userAggRes
	if err := cur.All(l.ctx, &usersAgg); err != nil {
		l.Errorw("Decode active users failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to decode active users")
	}

	// 4. 填充用户统计结果（简化：不调?User RPC，仅填充 userID ?count?
	activeUsers := make([]*msg.ActiveUser, 0, len(usersAgg))
	for _, u := range usersAgg {
		activeUsers = append(activeUsers, &msg.ActiveUser{
			User: &sdkws.UserInfo{
				UserID: u.ID,
			},
			Count: u.Count,
		})
	}

	return &msg.GetActiveUserResp{
		MsgCount:  totalMsgCount,
		UserCount: userCount,
		DateCount: dateCount,
		Users:     activeUsers,
	}, nil
}
