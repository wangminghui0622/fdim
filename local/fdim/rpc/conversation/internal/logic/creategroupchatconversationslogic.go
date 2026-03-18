package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/model"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupChatConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupChatConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupChatConversationsLogic {
	return &CreateGroupChatConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupChatConversationsLogic) CreateGroupChatConversations(req *conversation.CreateGroupChatConversationsReq) (*conversation.CreateGroupChatConversationsResp, error) {
	resp := &conversation.CreateGroupChatConversationsResp{}

	// 构建 conversationID（与 conversationutil.GenGroupConversationID 一致）
	conversationID := "sg_" + req.GroupID

	// 为每个用户创建群聊会话
	conversations := make([]*model.Conversation, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		conv := &model.Conversation{
			ConversationID:   conversationID,
			ConversationType: 2, // ReadGroupChatType
			GroupID:          req.GroupID,
			OwnerUserID:      userID,
			CreateTime:       time.Now(),
		}
		conversations = append(conversations, conv)
	}

	err := l.svcCtx.ConversationDB.Create(l.ctx, conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to create group chat conversations: %w", err)
	}

	// 初始化各用户该会话的 maxSeq（通常为 0），通知 Msg 模块
	if l.svcCtx.MsgClient != nil && len(req.UserIDs) > 0 {
		_, err = l.svcCtx.MsgClient.SetUserConversationMaxSeq(l.ctx, &msg.SetUserConversationMaxSeqReq{
			ConversationID: conversationID,
			OwnerUserID:    req.UserIDs,
			MaxSeq:         0,
		})
		if err != nil {
			l.Errorf("failed to init group conversation max seq: %v", err)
		}
	}

	return resp, nil
}
