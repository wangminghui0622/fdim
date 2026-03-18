package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationLogic {
	return &UpdateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConversationLogic) UpdateConversation(req *conversation.UpdateConversationReq) (*conversation.UpdateConversationResp, error) {
	resp := &conversation.UpdateConversationResp{}

	// 构建更新字段
	args := make(map[string]interface{})
	if req.RecvMsgOpt != nil {
		args["recv_msg_opt"] = req.RecvMsgOpt.Value
	}
	if req.IsPinned != nil {
		args["is_pinned"] = req.IsPinned.Value
	}
	if req.AttachedInfo != nil {
		args["attached_info"] = req.AttachedInfo.Value
	}

	if len(args) == 0 {
		return resp, nil
	}

	// 更新会话
	_, err := l.svcCtx.ConversationDB.UpdateByMap(l.ctx, req.UserIDs, req.ConversationID, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	// 发送会话变更通知（与官方一致）
	if l.svcCtx.ConversationNotification != nil && req.ConversationID != "" {
		for _, userID := range req.UserIDs {
			l.svcCtx.ConversationNotification.ConversationChangeNotification(l.ctx, userID, []string{req.ConversationID})
		}
	}

	return resp, nil
}
