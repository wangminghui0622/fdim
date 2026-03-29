package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetActiveConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActiveConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveConversationLogic {
	return &GetActiveConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetActiveConversationLogic) GetActiveConversation(req *msg.GetActiveConversationReq) (*msg.GetActiveConversationResp, error) {
	if len(req.ConversationIDs) == 0 {
		return &msg.GetActiveConversationResp{}, nil
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}

	coll := l.svcCtx.MongoDB.GetCollection("stream_msg")
	if coll == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message collection not initialized")
	}

	// 聚合：按会话分组，取最近发送时间和最?seq
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"conversation_id": bson.M{"$in": req.ConversationIDs},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":      "$conversation_id",
			"lastTime": bson.M{"$max": "$send_time"},
			"maxSeq":   bson.M{"$max": "$seq"},
		}}},
		{{Key: "$sort", Value: bson.M{"lastTime": -1}}},
		{{Key: "$limit", Value: req.Limit}},
	}

	cur, err := coll.Aggregate(l.ctx, pipeline)
	if err != nil {
		l.Errorw("Aggregate GetActiveConversation failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to aggregate active conversations")
	}
	defer cur.Close(l.ctx)

	type aggResult struct {
		ID       string `bson:"_id"`
		LastTime int64  `bson:"lastTime"`
		MaxSeq   int64  `bson:"maxSeq"`
	}

	var results []aggResult
	if err := cur.All(l.ctx, &results); err != nil {
		l.Errorw("Decode Aggregate result failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to decode active conversations")
	}

	resp := &msg.GetActiveConversationResp{
		Conversations: make([]*msg.ActiveConversation, 0, len(results)),
	}
	for _, r := range results {
		resp.Conversations = append(resp.Conversations, &msg.ActiveConversation{
			ConversationID: r.ID,
			LastTime:       r.LastTime,
			MaxSeq:         r.MaxSeq,
		})
	}

	return resp, nil
}
