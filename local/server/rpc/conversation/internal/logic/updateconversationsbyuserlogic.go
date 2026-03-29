package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationsByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationsByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationsByUserLogic {
	return &UpdateConversationsByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConversationsByUserLogic) UpdateConversationsByUser(req *conversation.UpdateConversationsByUserReq) (*conversation.UpdateConversationsByUserResp, error) {
	resp := &conversation.UpdateConversationsByUserResp{}

	// 构建更新字段
	args := make(map[string]interface{})
	if req.Ex != nil {
		args["ex"] = req.Ex.Value
	}

	if len(args) == 0 {
		return resp, nil
	}

	// 更新用户的所有会?
	_, err := l.svcCtx.ConversationDB.UpdateUserConversations(l.ctx, req.UserID, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversations: %w", err)
	}

	return resp, nil
}
