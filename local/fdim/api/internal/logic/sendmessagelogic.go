package logic

import (
	"context"
	"encoding/json"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMsgReq) (resp *types.SendMsgResp, err error) {
	// 验证 SendMsg 字段
	if req.SendMsg == nil {
		return nil, errs.ErrArgs.WrapMsg("field \"sendMsg\" is required")
	}

	// 验证必需字段
	if req.SendMsg.ClientMsgID == "" {
		return nil, errs.ErrArgs.WrapMsg("clientMsgID is required in sendMsg")
	}
	if req.SendMsg.SendID == "" {
		return nil, errs.ErrArgs.WrapMsg("sendID is required in sendMsg")
	}

	if req.RecvID != "" {
		req.SendMsg.RecvID = req.RecvID
	}

	msgData := &sdkws.MsgData{
		SendID:           req.SendMsg.SendID,
		RecvID:           req.SendMsg.RecvID,
		GroupID:          req.SendMsg.GroupID,
		ClientMsgID:      req.SendMsg.ClientMsgID,
		ServerMsgID:      req.SendMsg.ServerMsgID,
		SenderPlatformID: req.SendMsg.SenderPlatformID,
		SenderNickname:   req.SendMsg.SenderNickname,
		SenderFaceURL:    req.SendMsg.SenderFaceURL,
		SessionType:      req.SendMsg.SessionType,
		MsgFrom:          req.SendMsg.MsgFrom,
		ContentType:      req.SendMsg.ContentType,
		Content:          []byte(req.SendMsg.Content),
		Seq:              req.SendMsg.Seq,
		SendTime:         req.SendMsg.SendTime,
		CreateTime:       req.SendMsg.CreateTime,
		Status:           req.SendMsg.Status,
		IsRead:           req.SendMsg.IsRead,
		Options:          req.SendMsg.Options,
		AtUserIDList:     req.SendMsg.AtUserIDList,
		AttachedInfo:     req.SendMsg.AttachedInfo,
		Ex:               req.SendMsg.Ex,
	}

	// 处理 OfflinePushInfo
	if req.SendMsg.OfflinePushInfo != nil {
		jsonData, err := json.Marshal(req.SendMsg.OfflinePushInfo)
		if err == nil {
			var offlinePush sdkws.OfflinePushInfo
			if err := json.Unmarshal(jsonData, &offlinePush); err == nil {
				msgData.OfflinePushInfo = &offlinePush
			}
		}
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
	resp = &types.SendMsgResp{
		ServerMsgID: rpcResp.ServerMsgID,
		ClientMsgID: rpcResp.ClientMsgID,
		SendTime:    rpcResp.SendTime,
		Modify:      make(map[string]interface{}),
	}

	return resp, nil
}
