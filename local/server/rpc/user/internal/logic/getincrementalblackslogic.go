package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetIncrementalBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalBlacksLogic {
	return &GetIncrementalBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalBlacksLogic) GetIncrementalBlacks(req *user.GetIncrementalBlacksReq) (*user.GetIncrementalBlacksResp, error) {
	// TODO: 实现获取增量黑名单的逻辑
	return nil, fmt.Errorf("GetIncrementalBlacks not fully implemented yet")
}
