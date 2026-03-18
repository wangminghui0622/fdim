package logic

import (
	"context"
	"fmt"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
)

type SearchMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMsgLogic {
	return &SearchMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchMsgLogic) SearchMsg(req *types.SearchMsgReq) (resp *types.SearchMsgResp, err error) {
	// 转换请求参数
	var contentType int32
	if len(req.ContentType) > 0 {
		contentType = req.ContentType[0]
	}

	sendTimeStr := ""
	if req.SendTime > 0 {
		sendTimeStr = fmt.Sprintf("%d", req.SendTime)
	}

	rpcReq := &msg.SearchMessageReq{
		SendID:      req.SendID,
		RecvID:      req.RecvID,
		ContentType: contentType,
		SendTime:    sendTimeStr,
		SessionType: req.SessionType,
	}

	if req.Pagination.PageNumber > 0 || req.Pagination.ShowNumber > 0 {
		rpcReq.Pagination = &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		}
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.SearchMessage(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var msgs []interface{}
	for _, chatLog := range rpcResp.ChatLogs {
		msgs = append(msgs, chatLog)
	}

	resp = &types.SearchMsgResp{
		Total: int64(rpcResp.ChatLogsNum),
		Msgs:  msgs,
	}

	return resp, nil
}
