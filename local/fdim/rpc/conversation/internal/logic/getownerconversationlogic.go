package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOwnerConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOwnerConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOwnerConversationLogic {
	return &GetOwnerConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOwnerConversationLogic) GetOwnerConversation(req *conversation.GetOwnerConversationReq) (*conversation.GetOwnerConversationResp, error) {
	resp := &conversation.GetOwnerConversationResp{}

	// 获取分页参数
	offset := int64(0)
	limit := int64(20) // 默认值
	if req.Pagination != nil {
		offset = int64((req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber)
		limit = int64(req.Pagination.ShowNumber)
	}

	// 分页获取会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.PageConversationIDs(l.ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to page conversation IDs: %w", err)
	}

	// 获取总数
	total, err := l.svcCtx.ConversationDB.GetAllConversationIDsNumber(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total conversation IDs: %w", err)
	}

	// 根据会话ID获取会话
	conversations, err := l.svcCtx.ConversationDB.GetConversationsByConversationID(l.ctx, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// 转换为 protobuf 格式
	resp.Conversations = make([]*conversation.Conversation, 0, len(conversations))
	for _, conv := range conversations {
		resp.Conversations = append(resp.Conversations, &conversation.Conversation{
			OwnerUserID:      conv.OwnerUserID,
			ConversationID:   conv.ConversationID,
			ConversationType: conv.ConversationType,
			UserID:           conv.UserID,
			GroupID:          conv.GroupID,
			RecvMsgOpt:       conv.RecvMsgOpt,
			IsPinned:         conv.IsPinned,
			IsPrivateChat:    conv.IsPrivateChat,
			BurnDuration:     conv.BurnDuration,
			GroupAtType:      conv.GroupAtType,
			AttachedInfo:     conv.AttachedInfo,
			Ex:               conv.Ex,
			MaxSeq:           conv.MaxSeq,
			MinSeq:           conv.MinSeq,
			IsMsgDestruct:    conv.IsMsgDestruct,
			MsgDestructTime:  conv.MsgDestructTime,
		})
	}

	resp.Total = total
	return resp, nil
}
