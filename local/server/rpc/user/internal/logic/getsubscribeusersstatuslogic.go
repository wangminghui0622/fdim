package logic

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscribeUsersStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSubscribeUsersStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeUsersStatusLogic {
	return &GetSubscribeUsersStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSubscribeUsersStatusLogic) GetSubscribeUsersStatus(req *user.GetSubscribeUsersStatusReq) (*user.GetSubscribeUsersStatusResp, error) {
	// TODO: 实现获取订阅用户状态的逻辑
	// 当前返回空实现，待后续完善
	return &user.GetSubscribeUsersStatusResp{}, nil
}
