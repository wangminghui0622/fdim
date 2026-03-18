package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetNotNotifyConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNotNotifyConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotNotifyConversationIDsLogic {
	return &GetNotNotifyConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotNotifyConversationIDsLogic) GetNotNotifyConversationIDs(req *types.GetNotNotifyConversationIDsReq) (resp *types.GetNotNotifyConversationIDsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetNotNotifyConversationIDsReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetNotNotifyConversationIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetNotNotifyConversationIDsResp{
		ConversationIDs: rpcResp.ConversationIDs,
	}

	return resp, nil
}
