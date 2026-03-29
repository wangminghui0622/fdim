package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllConversationsLogic {
	return &GetAllConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllConversationsLogic) GetAllConversations(req *conversation.GetAllConversationsReq) (*conversation.GetAllConversationsResp, error) {
	resp := &conversation.GetAllConversationsResp{}

	// 查找用户的所有会?
	conversations, err := l.svcCtx.ConversationDB.FindUserIDAllConversations(l.ctx, req.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations: %w", err)
	}

	// 转换?protobuf 格式
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

	return resp, nil
}
