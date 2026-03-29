package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationOfflinePushUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationOfflinePushUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationOfflinePushUserIDsLogic {
	return &GetConversationOfflinePushUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationOfflinePushUserIDsLogic) GetConversationOfflinePushUserIDs(req *conversation.GetConversationOfflinePushUserIDsReq) (*conversation.GetConversationOfflinePushUserIDsResp, error) {
	resp := &conversation.GetConversationOfflinePushUserIDsResp{}

	if req.ConversationID == "" {
		return nil, fmt.Errorf("conversationID is empty")
	}
	if len(req.UserIDs) == 0 {
		return resp, nil
	}

	// 获取不接收消息的用户ID列表
	notReceiveUserIDs, err := l.svcCtx.ConversationDB.GetConversationNotReceiveMessageUserIDs(l.ctx, req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get not receive message user IDs: %w", err)
	}

	if len(notReceiveUserIDs) == 0 {
		resp.UserIDs = req.UserIDs
		return resp, nil
	}

	// 构建不接收消息的用户ID集合
	notReceiveSet := make(map[string]bool)
	for _, userID := range notReceiveUserIDs {
		notReceiveSet[userID] = true
	}

	// ?req.UserIDs 中排除不接收消息的用?
	var resultUserIDs []string
	for _, userID := range req.UserIDs {
		if !notReceiveSet[userID] {
			resultUserIDs = append(resultUserIDs, userID)
		}
	}

	resp.UserIDs = resultUserIDs
	return resp, nil
}
