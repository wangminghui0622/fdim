package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SortQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSortQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortQueryLogic {
	return &SortQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SortQueryLogic) SortQuery(req *user.SortQueryReq) (*user.SortQueryResp, error) {
	// TODO: 实现排序查询逻辑
	return nil, fmt.Errorf("SortQuery not implemented yet")
}
