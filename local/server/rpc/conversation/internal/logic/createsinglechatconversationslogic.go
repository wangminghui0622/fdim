package logic

import (
	"context"
	"time"

	"fdim/pkg/model"
	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSingleChatConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSingleChatConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSingleChatConversationsLogic {
	return &CreateSingleChatConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSingleChatConversationsLogic) CreateSingleChatConversations(req *conversation.CreateSingleChatConversationsReq) (*conversation.CreateSingleChatConversationsResp, error) {
	resp := &conversation.CreateSingleChatConversationsResp{}
	now := time.Now()

	// 官方行为：为双方各创建一条会话记录（OwnerUserID 不同，UserID 指向对方?
	convs := []*model.Conversation{
		{
			ConversationID:   req.ConversationID,
			ConversationType: req.ConversationType,
			OwnerUserID:      req.SendID,
			UserID:           req.RecvID,
			CreateTime:       now,
		},
		{
			ConversationID:   req.ConversationID,
			ConversationType: req.ConversationType,
			OwnerUserID:      req.RecvID,
			UserID:           req.SendID,
			CreateTime:       now,
		},
	}

	for _, conv := range convs {
		if err := l.svcCtx.ConversationDB.Create(l.ctx, []*model.Conversation{conv}); err != nil {
			l.Infow("conversation may already exist, ignoring create error",
				logx.Field("conversationID", conv.ConversationID),
				logx.Field("ownerUserID", conv.OwnerUserID),
				logx.Field("error", err))
		}
	}

	return resp, nil
}
