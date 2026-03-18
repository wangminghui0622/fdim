package logic

import (
	"context"
	"fmt"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
)

type DeleteLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogsLogic {
	return &DeleteLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogsLogic) DeleteLogs(req *third.DeleteLogsReq) (*third.DeleteLogsResp, error) {
	if l.svcCtx.MongoDB == nil {
		return nil, fmt.Errorf("mongo not initialized")
	}
	if len(req.LogIDs) == 0 {
		return &third.DeleteLogsResp{}, nil
	}

	coll := l.svcCtx.MongoDB.GetCollection("third_logs")
	filter := bson.M{"_id": bson.M{"$in": req.LogIDs}}

	if _, err := coll.DeleteMany(l.ctx, filter); err != nil {
		return nil, fmt.Errorf("failed to delete logs: %w", err)
	}

	return &third.DeleteLogsResp{}, nil
}
