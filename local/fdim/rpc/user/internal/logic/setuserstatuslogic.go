package logic

import (
	"context"
	"fmt"

	"fdim/protocol/sdkws"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserStatusLogic {
	return &SetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserStatusLogic) SetUserStatus(req *user.SetUserStatusReq) (*user.SetUserStatusResp, error) {
	resp := &user.SetUserStatusResp{}

	// ������֤
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// ����״̬�������߻�����
	var online []int32
	var offline []int32

	// ״̬������1=Online, 0=Offline
	if req.Status == 1 { // Online
		online = []int32{req.PlatformID}
	} else { // Offline
		offline = []int32{req.PlatformID}
	}

	// ���� Redis �е�����״̬
	if err := l.svcCtx.UserCache.SetUserOnlineStatus(l.ctx, req.UserID, online, offline); err != nil {
		return nil, err
	}

	// 发送用户状态变更通知（与官方一致）
	if l.svcCtx.UserNotification != nil {
		tips := &sdkws.UserStatusChangeTips{
			FromUserID: req.UserID,
			ToUserID:   req.UserID,
			Status:     req.Status,
			PlatformID: req.PlatformID,
		}
		l.svcCtx.UserNotification.UserStatusChangeNotification(l.ctx, tips)
	}

	return resp, nil
}
