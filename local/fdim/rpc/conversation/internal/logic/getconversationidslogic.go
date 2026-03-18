package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationIDsLogic {
	return &GetConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationIDsLogic) GetConversationIDs(req *conversation.GetConversationIDsReq) (*conversation.GetConversationIDsResp, error) {
	resp := &conversation.GetConversationIDsResp{}

	// 查找用户的所有会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.FindUserIDAllConversationID(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation IDs: %w", err)
	}

	resp.ConversationIDs = conversationIDs
	return resp, nil
}
