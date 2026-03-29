package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationsLogic {
	return &DeleteConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteConversationsLogic) DeleteConversations(req *conversation.DeleteConversationsReq) (*conversation.DeleteConversationsResp, error) {
	resp := &conversation.DeleteConversationsResp{}

	// 删除用户的会?
	err := l.svcCtx.ConversationDB.DeleteUsersConversations(l.ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to delete conversations: %w", err)
	}

	// 发送会话变更通知（与官方一致）
	if l.svcCtx.ConversationNotification != nil && len(req.ConversationIDs) > 0 {
		l.svcCtx.ConversationNotification.ConversationChangeNotification(l.ctx, req.OwnerUserID, req.ConversationIDs)
	}

	return resp, nil
}
