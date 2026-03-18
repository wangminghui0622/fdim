package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"fdim/Infrastructure_service/msggateway/internal/svc"
	"fdim/Infrastructure_service/msggateway/internal/ws"
	"fdim/protocol/msggateway"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineBatchPushOneMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOnlineBatchPushOneMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineBatchPushOneMsgLogic {
	return &OnlineBatchPushOneMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OnlineBatchPushOneMsgLogic) OnlineBatchPushOneMsg(req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	resp := &msggateway.OnlineBatchPushOneMsgResp{}

	if req.MsgData == nil {
		return resp, nil
	}

	if len(req.PushToUserIDs) == 0 {
		return resp, nil
	}

	// Serialize MsgData as JSON with Content as UTF-8 string (not base64)
	msgDataBytes, err := ws.MsgDataToJSON(req.MsgData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MsgData: %w", err)
	}

	// Wrap in official WsResp frame with WSPushMsg (2001)
	wsResp := ws.WsResp{
		ReqIdentifier: ws.WSPushMsg,
		ErrCode:       0,
		Data:          msgDataBytes,
	}
	msgBytes, err := json.Marshal(wsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal WsResp: %w", err)
	}

	// 批量推送消息到在线用户
	results := make([]*msggateway.SingleMsgToUserResults, 0, len(req.PushToUserIDs))
	for _, userID := range req.PushToUserIDs {
		result := &msggateway.SingleMsgToUserResults{
			UserID: userID,
			Resp:   make([]*msggateway.SingleMsgToUserPlatform, 0),
		}

		// 获取用户的所有连接
		sessions := l.svcCtx.WsServer.GetUserSessions(userID)
		if len(sessions) == 0 {
			// 用户不在线
			results = append(results, result)
			continue
		}

		// 向每个平台推送
		onlinePush := false
		for platformIDStr, session := range sessions {
			// 将 platformIDStr 转换为 int32
			platformID, err := strconv.ParseInt(platformIDStr, 10, 32)
			if err != nil {
				logx.Errorf("Failed to parse platformID %s: %v", platformIDStr, err)
				platformID = 0
			}

			platform := &msggateway.SingleMsgToUserPlatform{
				RecvPlatFormID: int32(platformID),
				ResultCode:     0,
			}

			// 推送消息
			if err := session.Write(msgBytes); err != nil {
				logx.Errorf("Failed to push to user %s platform %s: %v", userID, platformIDStr, err)
				platform.ResultCode = 1 // 推送失败
			} else {
				onlinePush = true
			}

			result.Resp = append(result.Resp, platform)
		}

		result.OnlinePush = onlinePush
		results = append(results, result)
	}

	resp.SinglePushResult = results
	return resp, nil
}
