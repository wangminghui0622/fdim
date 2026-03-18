package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type SetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationsLogic {
	return &SetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetConversationsLogic) SetConversations(req *types.SetConversationsReq) (resp *types.SetConversationsResp, err error) {
	// 转换请求参数
	var userIDs []string
	var conversationReq *conversation.ConversationReq

	if len(req.Conversations) > 0 {
		if conv, ok := req.Conversations[0].(map[string]interface{}); ok {
			conversationReq = &conversation.ConversationReq{}
			if conversationID, ok := conv["conversationID"].(string); ok {
				conversationReq.ConversationID = conversationID
			}
			if userID, ok := conv["userID"].(string); ok {
				conversationReq.UserID = userID
			}
			if groupID, ok := conv["groupID"].(string); ok {
				conversationReq.GroupID = groupID
			}
			if ownerUserID, ok := conv["ownerUserID"].(string); ok {
				userIDs = append(userIDs, ownerUserID)
			}
		}
	}

	rpcReq := &conversation.SetConversationsReq{
		UserIDs:      userIDs,
		Conversation: conversationReq,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ConversationClient.SetConversations(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SetConversationsResp{
	}

	return resp, nil
}
