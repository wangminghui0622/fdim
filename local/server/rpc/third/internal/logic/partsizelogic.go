package logic

import (
	"context"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PartSizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPartSizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PartSizeLogic {
	return &PartSizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PartSizeLogic) PartSize(req *third.PartSizeReq) (*third.PartSizeResp, error) {
	// 根据文件总大小，选择一个合适的分片 size，尽量让分片数量在 1 ~ 10000 之间
	const (
		minPartSize = int64(5 * 1024 * 1024)   // 5MB
		maxPartSize = int64(100 * 1024 * 1024) // 100MB
		maxNumSize  = int64(10000)
	)

	size := req.Size
	if size <= 0 {
		// 默认分片大小
		return &third.PartSizeResp{Size: minPartSize}, nil
	}

	// 目标：partSize = ceil(size / maxNumSize)，并限制在 [minPartSize, maxPartSize] 范围内
	partSize := size / maxNumSize
	if size%maxNumSize != 0 {
		partSize++
	}
	if partSize < minPartSize {
		partSize = minPartSize
	}
	if partSize > maxPartSize {
		partSize = maxPartSize
	}

	return &third.PartSizeResp{Size: partSize}, nil
}
