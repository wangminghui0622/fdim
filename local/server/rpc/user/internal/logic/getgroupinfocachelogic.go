package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoCacheLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoCacheLogic {
	return &GetGroupInfoCacheLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupInfoCacheLogic) GetGroupInfoCache(req *user.GetGroupInfoCacheReq) (*user.GetGroupInfoCacheResp, error) {
	// TODO: ʵֻȡȺϢ߼
	return nil, fmt.Errorf("GetGroupInfoCache not fully implemented yet")
}
