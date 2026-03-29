package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetIncrementalFriendsApplyToLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalFriendsApplyToLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalFriendsApplyToLogic {
	return &GetIncrementalFriendsApplyToLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalFriendsApplyToLogic) GetIncrementalFriendsApplyTo(req *user.GetIncrementalFriendsApplyToReq) (*user.GetIncrementalFriendsApplyToResp, error) {
	// TODO: ʵֻȡ루յ߼
	return nil, fmt.Errorf("GetIncrementalFriendsApplyTo not fully implemented yet")
}
