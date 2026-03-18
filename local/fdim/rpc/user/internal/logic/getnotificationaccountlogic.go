package logic

import (
	"context"
	"fmt"

	"fdim/pkg/constant"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotificationAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationAccountLogic {
	return &GetNotificationAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNotificationAccountLogic) GetNotificationAccount(req *user.GetNotificationAccountReq) (*user.GetNotificationAccountResp, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// 获取用户信息
	userInfo, err := l.svcCtx.UserDB.GetUserByID(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 检查用户是否是通知账户
	if userInfo.AppMangerLevel >= constant.AppAdmin {
		return &user.GetNotificationAccountResp{
			Account: &user.NotificationAccountInfo{
				UserID:         userInfo.UserID,
				FaceURL:        userInfo.FaceURL,
				NickName:       userInfo.Nickname,
				AppMangerLevel: userInfo.AppMangerLevel,
			},
		}, nil
	}

	return nil, fmt.Errorf("notification messages cannot be sent for this ID")
}
