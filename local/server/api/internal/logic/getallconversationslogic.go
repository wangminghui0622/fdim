package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetAllConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllConversationsLogic {
	return &GetAllConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllConversationsLogic) GetAllConversations(req *types.GetAllConversationsReq) (resp *types.GetAllConversationsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetAllConversationsReq{
		OwnerUserID: req.OwnerUserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetAllConversations(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var conversations []interface{}
	for _, c := range rpcResp.Conversations {
		conversations = append(conversations, c)
	}

	resp = &types.GetAllConversationsResp{
		Conversations: conversations,
	}

	return resp, nil
}
