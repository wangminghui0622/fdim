package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationNotReceiveMessageUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationNotReceiveMessageUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationNotReceiveMessageUserIDsLogic {
	return &GetConversationNotReceiveMessageUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationNotReceiveMessageUserIDsLogic) GetConversationNotReceiveMessageUserIDs(req *conversation.GetConversationNotReceiveMessageUserIDsReq) (*conversation.GetConversationNotReceiveMessageUserIDsResp, error) {
	resp := &conversation.GetConversationNotReceiveMessageUserIDsResp{}

	// 获取不接收消息的用户ID列表
	userIDs, err := l.svcCtx.ConversationDB.GetConversationNotReceiveMessageUserIDs(l.ctx, req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get not receive message user IDs: %w", err)
	}

	resp.UserIDs = userIDs
	return resp, nil
}
