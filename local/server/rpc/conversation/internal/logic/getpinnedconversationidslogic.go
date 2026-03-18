package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPinnedConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPinnedConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPinnedConversationIDsLogic {
	return &GetPinnedConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPinnedConversationIDsLogic) GetPinnedConversationIDs(req *conversation.GetPinnedConversationIDsReq) (*conversation.GetPinnedConversationIDsResp, error) {
	resp := &conversation.GetPinnedConversationIDsResp{}

	// 查找用户的所有置顶会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.FindUserIDAllPinnedConversationID(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pinned conversation IDs: %w", err)
	}

	resp.ConversationIDs = conversationIDs
	return resp, nil
}
