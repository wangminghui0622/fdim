package logic

import (
	"context"
	"fdim/Infrastructure_service/msggateway/internal/svc"
	"strconv"

	"fdim/protocol/msggateway"
	"github.com/zeromicro/go-zero/core/logx"
)

type KickUserOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickUserOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickUserOfflineLogic {
	return &KickUserOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickUserOfflineLogic) KickUserOffline(req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	l.Infof("[Kick User Offline] Request received - platformID=%d, userCount=%d", req.PlatformID, len(req.KickUserIDList))
	
	// 从 WebSocket 服务器踢用户下线
	for _, userID := range req.KickUserIDList {
		// 如果指定了 platformID，只踢指定平台的连接
		if req.PlatformID != 0 {
			// 修复：正确转换 platformID 为字符串
			platformIDStr := strconv.Itoa(int(req.PlatformID))
			l.Infof("[Kick User Offline] Kicking user from specific platform - userID=%s, platformID=%s", userID, platformIDStr)
			if err := l.svcCtx.WsServer.KickUserOffline(userID, platformIDStr); err != nil {
				l.Errorf("[Kick User Offline] Failed to kick user from platform - userID=%s, platformID=%s, error=%v", userID, platformIDStr, err)
			} else {
				l.Infof("[Kick User Offline] Successfully kicked user - userID=%s, platformID=%s", userID, platformIDStr)
			}
		} else {
			// 如果没有指定 platformID，踢所有平台的连接
			l.Infof("[Kick User Offline] Kicking user from all platforms - userID=%s", userID)
			sessions := l.svcCtx.WsServer.GetUserSessions(userID)
			kickedCount := 0
			for platformIDStr := range sessions {
				if err := l.svcCtx.WsServer.KickUserOffline(userID, platformIDStr); err != nil {
					l.Errorf("[Kick User Offline] Failed to kick user from platform - userID=%s, platformID=%s, error=%v", userID, platformIDStr, err)
				} else {
					kickedCount++
				}
			}
			l.Infof("[Kick User Offline] Kicked user from all platforms - userID=%s, kickedCount=%d", userID, kickedCount)
		}
	}

	l.Infof("[Kick User Offline] Request completed - platformID=%d, userCount=%d", req.PlatformID, len(req.KickUserIDList))
	return &msggateway.KickUserOfflineResp{}, nil
}
