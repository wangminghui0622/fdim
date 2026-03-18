package logic

import (
	"context"
	"fmt"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetAppBadgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetAppBadgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAppBadgeLogic {
	return &SetAppBadgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetAppBadgeLogic) SetAppBadge(req *third.SetAppBadgeReq) (*third.SetAppBadgeResp, error) {
	// 1. 参数检查
	if err := req.Check(); err != nil {
		return nil, fmt.Errorf("invalid SetAppBadgeReq: %w", err)
	}
	if l.svcCtx.Redis == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	// 2. 将应用未读数写入 Redis，方便后续扩展（如 iOS 角标、统计等）
	key := fmt.Sprintf("app_badge:%s", req.UserID)
	if err := l.svcCtx.Redis.Set(l.ctx, key, req.AppUnreadCount, 0).Err(); err != nil {
		l.Errorf("failed to set app badge in redis, key=%s, err=%v", key, err)
		return nil, fmt.Errorf("failed to save app badge: %w", err)
	}

	l.Infof("SetAppBadge success, userID=%s, appUnreadCount=%d", req.UserID, req.AppUnreadCount)
	return &third.SetAppBadgeResp{}, nil
}
