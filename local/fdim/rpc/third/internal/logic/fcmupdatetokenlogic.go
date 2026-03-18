package logic

import (
	"context"
	"fmt"

	"fdim/protocol/third"
	"fdim/rpc/third/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FcmUpdateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFcmUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FcmUpdateTokenLogic {
	return &FcmUpdateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FcmUpdateTokenLogic) FcmUpdateToken(req *third.FcmUpdateTokenReq) (*third.FcmUpdateTokenResp, error) {
	// 1. 参数检查（使用 proto 自带的 Check 方法）
	if err := req.Check(); err != nil {
		return nil, fmt.Errorf("invalid FcmUpdateTokenReq: %w", err)
	}
	if l.svcCtx.Redis == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	// 2. 将 FCM token 写入 Redis，供 push 服务离线推送使用
	// key 规则与 push 服务中 DelUserPushToken 保持一致：fcm_token:{account}:{platformID}
	key := fmt.Sprintf("fcm_token:%s:%d", req.Account, req.PlatformID)
	if err := l.svcCtx.Redis.Set(l.ctx, key, req.FcmToken, 0).Err(); err != nil {
		l.Errorf("failed to set FCM token in redis, key=%s, err=%v", key, err)
		return nil, fmt.Errorf("failed to save FCM token: %w", err)
	}

	l.Infof("FcmUpdateToken success, account=%s, platformID=%d", req.Account, req.PlatformID)
	return &third.FcmUpdateTokenResp{}, nil
}
