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

type OnlinePushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOnlinePushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlinePushMsgLogic {
	return &OnlinePushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OnlinePushMsgLogic) OnlinePushMsg(req *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	resp := &msggateway.OnlinePushMsgResp{
		Resp: make([]*msggateway.SingleMsgToUserPlatform, 0),
	}

	logx.Infof("[OnlinePushMsg] Received push request for user: %s, clientMsgID: %s",
		req.PushToUserID, req.MsgData.ClientMsgID)

	if req.MsgData == nil {
		return nil, fmt.Errorf("msgData is nil")
	}

	if req.PushToUserID == "" {
		logx.Info("[OnlinePushMsg] PushToUserID is empty, skipping")
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
	respBytes, err := json.Marshal(wsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal WsResp: %w", err)
	}

	logx.Infof("[OnlinePushMsg] Broadcasting to user %s, message size: %d bytes", req.PushToUserID, len(respBytes))

	sessions := l.svcCtx.WsServer.GetUserSessions(req.PushToUserID)
	if len(sessions) == 0 {
		logx.Infof("[OnlinePushMsg] User %s not online, returning empty response", req.PushToUserID)
		return resp, nil
	}

	onlinePush := false
	for platformIDStr, session := range sessions {
		platformID, err := strconv.ParseInt(platformIDStr, 10, 32)
		if err != nil {
			logx.Errorf("[OnlinePushMsg] Failed to parse platformID %s: %v", platformIDStr, err)
			platformID = 0
		}

		platform := &msggateway.SingleMsgToUserPlatform{
			RecvID:         req.PushToUserID,
			RecvPlatFormID: int32(platformID),
			ResultCode:     0,
		}

		if err := session.Write(respBytes); err != nil {
			logx.Errorf("[OnlinePushMsg] Failed to push message to user %s platform %s: %v", req.PushToUserID, platformIDStr, err)
			platform.ResultCode = 1
		} else {
			onlinePush = true
			logx.Infof("[OnlinePushMsg] Successfully wrote message to user %s platform %s", req.PushToUserID, platformIDStr)
		}

		resp.Resp = append(resp.Resp, platform)
	}

	if !onlinePush {
		logx.Infof("[OnlinePushMsg] Broadcast to user %s failed on all %d sessions", req.PushToUserID, len(sessions))
	}

	logx.Infof("[OnlinePushMsg] Successfully broadcasted message to user: %s", req.PushToUserID)
	return resp, nil
}
