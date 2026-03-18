package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/model"
	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationLogic {
	return &SetConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetConversationLogic) SetConversation(req *conversation.SetConversationReq) (*conversation.SetConversationResp, error) {
	resp := &conversation.SetConversationResp{}

	if req.Conversation == nil {
		return nil, fmt.Errorf("conversation is nil")
	}

	// 转换为数据库模型
	conv := &model.Conversation{
		OwnerUserID:      req.Conversation.OwnerUserID,
		ConversationID:   req.Conversation.ConversationID,
		ConversationType: req.Conversation.ConversationType,
		UserID:           req.Conversation.UserID,
		GroupID:          req.Conversation.GroupID,
		RecvMsgOpt:       req.Conversation.RecvMsgOpt,
		IsPinned:         req.Conversation.IsPinned,
		IsPrivateChat:    req.Conversation.IsPrivateChat,
		BurnDuration:     req.Conversation.BurnDuration,
		GroupAtType:      req.Conversation.GroupAtType,
		AttachedInfo:     req.Conversation.AttachedInfo,
		Ex:               req.Conversation.Ex,
		MaxSeq:           req.Conversation.MaxSeq,
		MinSeq:           req.Conversation.MinSeq,
		IsMsgDestruct:    req.Conversation.IsMsgDestruct,
		MsgDestructTime:  req.Conversation.MsgDestructTime,
		CreateTime:       time.Now(),
	}

	// 更新或创建会话
	err := l.svcCtx.ConversationDB.Update(l.ctx, conv)
	if err != nil {
		// 如果不存在，则创建
		err = l.svcCtx.ConversationDB.Create(l.ctx, []*model.Conversation{conv})
		if err != nil {
			return nil, fmt.Errorf("failed to set conversation: %w", err)
		}
	}

	// 发送会话变更通知（与官方一致）
	if l.svcCtx.ConversationNotification != nil && req.Conversation.ConversationID != "" {
		l.svcCtx.ConversationNotification.ConversationChangeNotification(l.ctx, req.Conversation.OwnerUserID, []string{req.Conversation.ConversationID})
	}

	return resp, nil
}
