package logic

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeOrCancelUsersStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubscribeOrCancelUsersStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeOrCancelUsersStatusLogic {
	return &SubscribeOrCancelUsersStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubscribeOrCancelUsersStatusLogic) SubscribeOrCancelUsersStatus(req *user.SubscribeOrCancelUsersStatusReq) (*user.SubscribeOrCancelUsersStatusResp, error) {
	// TODO: 实现订阅/取消订阅用户状态的逻辑
	// 当前返回空实现
	return &user.SubscribeOrCancelUsersStatusResp{}, nil
}
