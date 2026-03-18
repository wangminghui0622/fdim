package logic

import (
	"context"
	"fmt"

	"fdim/Infrastructure_service/push/internal/svc"
	"fdim/pkg/cache"
	"fdim/protocol/push"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserPushTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserPushTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserPushTokenLogic {
	return &DelUserPushTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserPushTokenLogic) DelUserPushToken(req *push.DelUserPushTokenReq) (*push.DelUserPushTokenResp, error) {
	// 1. 验证参数
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// 2. 从 Redis 中删除用户的推送Token
	// FCM Token 的 key 格式: fcm_token:{userID}:{platformID}
	redisClient := cache.NewRedisClient(l.svcCtx.Config.Cache)
	key := fmt.Sprintf("fcm_token:%s:%d", req.UserID, req.PlatformID)

	if err := redisClient.Del(l.ctx, key).Err(); err != nil {
		l.Errorf("Failed to delete FCM token from Redis: %v", err)
		return nil, fmt.Errorf("failed to delete FCM token: %w", err)
	}

	l.Infof("Deleted FCM token for user: userID=%s, platformID=%d", req.UserID, req.PlatformID)

	// 3. 返回结果
	return &push.DelUserPushTokenResp{}, nil
}
