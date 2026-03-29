package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsByConversationIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsByConversationIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsByConversationIDLogic {
	return &GetConversationsByConversationIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationsByConversationIDLogic) GetConversationsByConversationID(req *conversation.GetConversationsByConversationIDReq) (*conversation.GetConversationsByConversationIDResp, error) {
	resp := &conversation.GetConversationsByConversationIDResp{}

	// 根据会话ID查找会话
	conversations, err := l.svcCtx.ConversationDB.GetConversationsByConversationID(l.ctx, req.ConversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
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
