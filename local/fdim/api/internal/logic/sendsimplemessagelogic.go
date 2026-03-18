package logic

import (
	"context"
	"encoding/json"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
)

type SendSimpleMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSimpleMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSimpleMessageLogic {
	return &SendSimpleMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSimpleMessageLogic) SendSimpleMessage(req *types.SendSimpleMsgReq) (resp *types.SendSimpleMsgResp, err error) {
	// 转换请求参数
	// 创建 MarkdownTextElem
	markdownContent := map[string]string{
		"content": req.Content,
	}
	contentJson, _ := json.Marshal(markdownContent)

	msgData := &sdkws.MsgData{
		SendID:      req.SendID,
		Content:     contentJson,
		ContentType: 1400, // MarkdownText
		SessionType: 1,    // SingleChatType
	}

	if req.OfflinePushInfo != nil {
		offlinePushJson, _ := json.Marshal(req.OfflinePushInfo)
		var offlinePush sdkws.OfflinePushInfo
		json.Unmarshal(offlinePushJson, &offlinePush)
		msgData.OfflinePushInfo = &offlinePush
	}

	if req.Ex != "" {
		msgData.Ex = req.Ex
	}

	rpcReq := &msg.SendSimpleMsgReq{
		MsgData: msgData,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.SendSimpleMsg(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SendSimpleMsgResp{
		ServerMsgID: rpcResp.ServerMsgID,
		ClientMsgID: rpcResp.ClientMsgID,
		SendTime:    rpcResp.SendTime,
	}

	return resp, nil
}
