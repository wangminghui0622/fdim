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

type SetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationsLogic {
	return &SetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetConversationsLogic) SetConversations(req *conversation.SetConversationsReq) (*conversation.SetConversationsResp, error) {
	resp := &conversation.SetConversationsResp{}

	if req.Conversation == nil {
		return nil, fmt.Errorf("conversation is nil")
	}

	// 构建更新字段
	args := make(map[string]interface{})
	if req.Conversation.RecvMsgOpt != nil {
		args["recv_msg_opt"] = req.Conversation.RecvMsgOpt.Value
	}
	if req.Conversation.IsPinned != nil {
		args["is_pinned"] = req.Conversation.IsPinned.Value
	}
	if req.Conversation.AttachedInfo != nil {
		args["attached_info"] = req.Conversation.AttachedInfo.Value
	}
	if req.Conversation.IsPrivateChat != nil {
		args["is_private_chat"] = req.Conversation.IsPrivateChat.Value
	}
	if req.Conversation.Ex != nil {
		args["ex"] = req.Conversation.Ex.Value
	}
	if req.Conversation.BurnDuration != nil {
		args["burn_duration"] = req.Conversation.BurnDuration.Value
	}
	if req.Conversation.MinSeq != nil {
		args["min_seq"] = req.Conversation.MinSeq.Value
	}
	if req.Conversation.MaxSeq != nil {
		args["max_seq"] = req.Conversation.MaxSeq.Value
	}
	if req.Conversation.GroupAtType != nil {
		args["group_at_type"] = req.Conversation.GroupAtType.Value
	}
	if req.Conversation.MsgDestructTime != nil {
		args["msg_destruct_time"] = req.Conversation.MsgDestructTime.Value
	}
	if req.Conversation.IsMsgDestruct != nil {
		args["is_msg_destruct"] = req.Conversation.IsMsgDestruct.Value
	}

	// 收集变更的会话ID，用于发送通知
	conversationIDs := []string{}
	if req.Conversation.ConversationID != "" {
		conversationIDs = append(conversationIDs, req.Conversation.ConversationID)
	}

	// 为每个用户更新或创建会话
	for _, userID := range req.UserIDs {
		// 检查会话是否存在
		existing, err := l.svcCtx.ConversationDB.Find(l.ctx, userID, []string{req.Conversation.ConversationID})
		if err != nil {
			l.Errorf("failed to find conversation: %v", err)
			continue
		}

		if len(existing) > 0 {
			// 更新现有会话
			if len(args) > 0 {
				_, err := l.svcCtx.ConversationDB.UpdateByMap(l.ctx, []string{userID}, req.Conversation.ConversationID, args)
				if err != nil {
					l.Errorf("failed to update conversation: %v", err)
				}
			}
		} else {
			// 创建新会话
			conv := &model.Conversation{
				ConversationID:   req.Conversation.ConversationID,
				ConversationType: req.Conversation.ConversationType,
				UserID:           req.Conversation.UserID,
				GroupID:          req.Conversation.GroupID,
				OwnerUserID:      userID,
				CreateTime:       time.Now(),
			}
			// 设置可选字段
			if req.Conversation.RecvMsgOpt != nil {
				conv.RecvMsgOpt = req.Conversation.RecvMsgOpt.Value
			}
			if req.Conversation.IsPinned != nil {
				conv.IsPinned = req.Conversation.IsPinned.Value
			}
			if req.Conversation.AttachedInfo != nil {
				conv.AttachedInfo = req.Conversation.AttachedInfo.Value
			}
			if req.Conversation.IsPrivateChat != nil {
				conv.IsPrivateChat = req.Conversation.IsPrivateChat.Value
			}
			if req.Conversation.Ex != nil {
				conv.Ex = req.Conversation.Ex.Value
			}

			err := l.svcCtx.ConversationDB.Create(l.ctx, []*model.Conversation{conv})
			if err != nil {
				l.Errorf("failed to create conversation: %v", err)
			}
		}
	}

	// 发送会话变更通知（与官方一致）
	if l.svcCtx.ConversationNotification != nil && len(conversationIDs) > 0 {
		for _, userID := range req.UserIDs {
			l.svcCtx.ConversationNotification.ConversationChangeNotification(l.ctx, userID, conversationIDs)
		}
	}

	return resp, nil
}
