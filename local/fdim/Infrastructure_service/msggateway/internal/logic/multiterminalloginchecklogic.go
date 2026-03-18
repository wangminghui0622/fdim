package logic

import (
	"context"
	"fdim/Infrastructure_service/msggateway/internal/svc"
	"strconv"

	"fdim/protocol/msggateway"
	"github.com/zeromicro/go-zero/core/logx"
)

type MultiTerminalLoginCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMultiTerminalLoginCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiTerminalLoginCheckLogic {
	return &MultiTerminalLoginCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MultiTerminalLoginCheckLogic) MultiTerminalLoginCheck(req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {
	resp := &msggateway.MultiTerminalLoginCheckResp{}

	// 检查用户是否已经在其他终端登录
	userID := req.UserID
	platformID := req.PlatformID

	l.Infof("[Multi-Terminal Check] Checking login - userID=%s, platformID=%d", userID, platformID)

	// 获取用户的所有会话
	sessions := l.svcCtx.WsServer.GetUserSessions(userID)

	// 如果用户已经在其他平台登录，需要踢掉旧连接
	if len(sessions) > 0 {
		// 修复：正确转换 platformID 为字符串
		platformIDStr := strconv.Itoa(int(platformID))
		kickedCount := 0
		for existingPlatformIDStr := range sessions {
			if existingPlatformIDStr != platformIDStr {
				// 踢掉旧连接
				l.Infof("[Multi-Terminal Check] Kicking old session - userID=%s, oldPlatform=%s, newPlatform=%s",
					userID, existingPlatformIDStr, platformIDStr)
				if err := l.svcCtx.WsServer.KickUserOffline(userID, existingPlatformIDStr); err != nil {
					l.Errorf("[Multi-Terminal Check] Failed to kick old session - userID=%s, platform=%s, error=%v",
						userID, existingPlatformIDStr, err)
				} else {
					kickedCount++
				}
			}
		}
		l.Infof("[Multi-Terminal Check] Completed - userID=%s, platformID=%d, kickedCount=%d", userID, platformID, kickedCount)
	}

	return resp, nil
}
