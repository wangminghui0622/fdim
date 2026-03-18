package logic

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationUserInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationUserInfoUpdateLogic {
	return &NotificationUserInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationUserInfoUpdateLogic) NotificationUserInfoUpdate(req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	// TODO: 实现通知用户信息更新的逻辑
	// 当前返回空实现
	return &user.NotificationUserInfoUpdateResp{}, nil
}
