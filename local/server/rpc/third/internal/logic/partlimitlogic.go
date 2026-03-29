package logic

import (
	"context"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PartLimitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPartLimitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PartLimitLogic {
	return &PartLimitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PartLimitLogic) PartLimit(req *third.PartLimitReq) (*third.PartLimitResp, error) {
	const (
		minPartSize = int64(5 * 1024 * 1024)   // 5MB
		maxPartSize = int64(100 * 1024 * 1024) // 100MB
		maxNumSize  = int32(10000)
	)

	return &third.PartLimitResp{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}, nil
}
