package logic

import (
	"context"
	"encoding/json"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
)

type SendBusinessNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendBusinessNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBusinessNotificationLogic {
	return &SendBusinessNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendBusinessNotificationLogic) SendBusinessNotification(req *types.SendBusinessNotificationReq) (resp *types.SendBusinessNotificationResp, err error) {
	// 转换请求参数
	var sessionType int32
	var recvID string
	if req.RecvUserID != "" {
		sessionType = 1 // SingleChatType
		recvID = req.RecvUserID
	} else {
		sessionType = 2 // ReadGroupChatType
		recvID = ""
	}

	notificationData := map[string]string{
		"key":  req.Key,
		"data": req.Data,
	}
	notificationJson, _ := json.Marshal(notificationData)

	notificationElem := &sdkws.NotificationElem{
		Detail: string(notificationJson),
	}
	notificationContent, _ := json.Marshal(notificationElem)

	msgData := &sdkws.MsgData{
		SendID:      req.SendUserID,
		RecvID:      recvID,
		GroupID:     req.RecvGroupID,
		Content:     notificationContent,
		ContentType: 1001, // BusinessNotification
		SessionType: sessionType,
	}

	rpcReq := &msg.SendMsgReq{
		MsgData: msgData,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.SendMsg(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SendBusinessNotificationResp{
		ServerMsgID: rpcResp.ServerMsgID,
		ClientMsgID: rpcResp.ClientMsgID,
		SendTime:    rpcResp.SendTime,
	}

	return resp, nil
}
