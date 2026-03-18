package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFullOwnerConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFullOwnerConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullOwnerConversationIDsLogic {
	return &GetFullOwnerConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFullOwnerConversationIDsLogic) GetFullOwnerConversationIDs(req *conversation.GetFullOwnerConversationIDsReq) (*conversation.GetFullOwnerConversationIDsResp, error) {
	resp := &conversation.GetFullOwnerConversationIDsResp{}

	// 查找用户的所有会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.FindUserIDAllConversationID(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation IDs: %w", err)
	}

	resp.ConversationIDs = conversationIDs
	return resp, nil
}
