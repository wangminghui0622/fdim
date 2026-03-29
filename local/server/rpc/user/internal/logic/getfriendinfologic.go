package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendInfoLogic {
	return &GetFriendInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendInfoLogic) GetFriendInfo(req *user.GetFriendInfoReq) (*user.GetFriendInfoResp, error) {
	// TODO: ʵֻȡϢ߼
	return nil, fmt.Errorf("GetFriendInfo not fully implemented yet")
}
