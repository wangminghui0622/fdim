package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationLogic {
	return &GetConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationLogic) GetConversation(req *conversation.GetConversationReq) (*conversation.GetConversationResp, error) {
	resp := &conversation.GetConversationResp{}

	// 获取单个会话
	conv, err := l.svcCtx.ConversationDB.Take(l.ctx, req.OwnerUserID, req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// 转换?protobuf 格式
	resp.Conversation = &conversation.Conversation{
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
	}

	return resp, nil
}
