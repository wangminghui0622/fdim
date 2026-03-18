package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotNotifyConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotNotifyConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotNotifyConversationIDsLogic {
	return &GetNotNotifyConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNotNotifyConversationIDsLogic) GetNotNotifyConversationIDs(req *conversation.GetNotNotifyConversationIDsReq) (*conversation.GetNotNotifyConversationIDsResp, error) {
	resp := &conversation.GetNotNotifyConversationIDsResp{}

	// 查找用户的所有不通知会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.FindUserIDAllNotNotifyConversationID(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find not notify conversation IDs: %w", err)
	}

	resp.ConversationIDs = conversationIDs
	return resp, nil
}
