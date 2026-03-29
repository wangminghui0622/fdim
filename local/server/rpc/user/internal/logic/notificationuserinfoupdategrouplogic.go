package logic

import (
	"context"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationUserInfoUpdateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotificationUserInfoUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationUserInfoUpdateGroupLogic {
	return &NotificationUserInfoUpdateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotificationUserInfoUpdateGroupLogic) NotificationUserInfoUpdate(req *user.NotificationUserInfoUpdateReq) (*user.NotificationUserInfoUpdateResp, error) {
	// TODO: ʵ֪ͨûϢµ߼
	// ǰؿʵ
	return &user.NotificationUserInfoUpdateResp{}, nil
}
