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
	// 分片限制这里直接采用与 open-im-server 接近的一组默认值，便于客户端统一处理
	// 一般 S3 兼容存储：最小分片 5MB，最大分片 100MB，最大分片数 10000
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
