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

type DeleteOutdatedDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOutdatedDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOutdatedDataLogic {
	return &DeleteOutdatedDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteOutdatedDataLogic) DeleteOutdatedData(req *third.DeleteOutdatedDataReq) (*third.DeleteOutdatedDataResp, error) {
	if err := req.Check(); err != nil {
		return nil, fmt.Errorf("invalid DeleteOutdatedDataReq: %w", err)
	}
	if l.svcCtx.MongoDB == nil {
		return nil, fmt.Errorf("mongo not initialized")
	}

	coll := l.svcCtx.MongoDB.GetCollection("third_data")

	filter := bson.M{
		"group": bson.M{"$in": req.ObjectGroup},
	}
	opts := options.Find().SetSort(bson.D{{Key: "update_time", Value: 1}}).SetLimit(int64(req.Limit))

	cursor, err := coll.Find(l.ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query outdated data: %w", err)
	}
	defer cursor.Close(l.ctx)

	var ids []interface{}
	for cursor.Next(l.ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode data: %w", err)
		}
		if id, ok := doc["_id"]; ok {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return &third.DeleteOutdatedDataResp{Count: 0}, nil
	}

	delFilter := bson.M{"_id": bson.M{"$in": ids}}
	res, err := coll.DeleteMany(l.ctx, delFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete outdated data: %w", err)
	}

	return &third.DeleteOutdatedDataResp{Count: int32(res.DeletedCount)}, nil
}
