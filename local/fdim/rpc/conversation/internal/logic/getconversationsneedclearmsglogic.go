package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsNeedClearMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsNeedClearMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsNeedClearMsgLogic {
	return &GetConversationsNeedClearMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationsNeedClearMsgLogic) GetConversationsNeedClearMsg(req *conversation.GetConversationsNeedClearMsgReq) (*conversation.GetConversationsNeedClearMsgResp, error) {
	resp := &conversation.GetConversationsNeedClearMsgResp{}

	// 获取需要销毁的会话
	conversations, err := l.svcCtx.ConversationDB.GetConversationIDsNeedDestruct(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations need clear msg: %w", err)
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

	return resp, nil
}
