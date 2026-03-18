package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationLogic {
	return &GetConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationLogic) GetConversation(req *types.GetConversationReq) (resp *types.GetConversationResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetConversationReq{
		OwnerUserID:    req.OwnerUserID,
		ConversationID: req.ConversationID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetConversation(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetConversationResp{
		Conversation: rpcResp.Conversation,
	}

	return resp, nil
}
