package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberCacheLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberCacheLogic {
	return &GetGroupMemberCacheLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberCacheLogic) GetGroupMemberCache(req *user.GetGroupMemberCacheReq) (*user.GetGroupMemberCacheResp, error) {
	// TODO: ʵֻȡȺԱ߼
	return nil, fmt.Errorf("GetGroupMemberCache not fully implemented yet")
}
