package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetConversationsReq{
		OwnerUserID:     req.OwnerUserID,
		ConversationIDs: req.ConversationIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetConversations(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var conversations []interface{}
	for _, c := range rpcResp.Conversations {
		conversations = append(conversations, c)
	}

	resp = &types.GetConversationsResp{
		Conversations: conversations,
	}

	return resp, nil
}
