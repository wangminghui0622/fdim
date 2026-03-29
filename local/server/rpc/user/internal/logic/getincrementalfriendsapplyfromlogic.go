package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetIncrementalFriendsApplyFromLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalFriendsApplyFromLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalFriendsApplyFromLogic {
	return &GetIncrementalFriendsApplyFromLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalFriendsApplyFromLogic) GetIncrementalFriendsApplyFrom(req *user.GetIncrementalFriendsApplyFromReq) (*user.GetIncrementalFriendsApplyFromResp, error) {
	// TODO: ʵֻȡ루߼
	return nil, fmt.Errorf("GetIncrementalFriendsApplyFrom not fully implemented yet")
}
