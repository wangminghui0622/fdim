package logic

import (
	"context"
	"fdim/Infrastructure_service/msggateway/internal/svc"
	"strconv"

	"fdim/protocol/msggateway"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersOnlineStatusLogic {
	return &GetUsersOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersOnlineStatusLogic) GetUsersOnlineStatus(req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	resp := &msggateway.GetUsersOnlineStatusResp{}

	if len(req.UserIDs) == 0 {
		return resp, nil
	}

	// ?WebSocket 服务器获取在线状?
	onlineUsers := l.svcCtx.WsServer.GetOnlineUsers(req.UserIDs)

	// 构建响应
	for _, userID := range req.UserIDs {
		if online, ok := onlineUsers[userID]; ok && online {
			// 获取该用户的所有连?
			sessions := l.svcCtx.WsServer.GetUserSessions(userID)
			successResult := &msggateway.GetUsersOnlineStatusResp_SuccessResult{
				UserID: userID,
				Status: 1, // 在线
			}

			// 为每个平台构建详?
			for platformIDStr, session := range sessions {
				platformIDInt, err := strconv.Atoi(platformIDStr)
				if err != nil {
					logx.Errorf("Failed to parse platformID: %v", err)
					continue
				}

				// ?session 中获取信?
				connID := session.RemoteAddr().String() // 使用 RemoteAddr 作为 connID
				isBackground := false                   // 默认?false，可以从 session 中获?
				detail := &msggateway.GetUsersOnlineStatusResp_SuccessDetail{
					PlatformID:   int32(platformIDInt),
					ConnID:       connID,
					IsBackground: isBackground,
				}
				successResult.DetailPlatformStatus = append(successResult.DetailPlatformStatus, detail)
			}

			resp.SuccessResult = append(resp.SuccessResult, successResult)
		} else {
			// 用户不在?
			failedDetail := &msggateway.GetUsersOnlineStatusResp_FailedDetail{
				UserID: userID,
			}
			resp.FailedResult = append(resp.FailedResult, failedDetail)
		}
	}

	return resp, nil
}
