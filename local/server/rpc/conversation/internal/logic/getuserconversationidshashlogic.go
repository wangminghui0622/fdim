package logic

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationIDsHashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConversationIDsHashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationIDsHashLogic {
	return &GetUserConversationIDsHashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConversationIDsHashLogic) GetUserConversationIDsHash(req *conversation.GetUserConversationIDsHashReq) (*conversation.GetUserConversationIDsHashResp, error) {
	resp := &conversation.GetUserConversationIDsHashResp{}

	// 获取用户的所有会话ID
	conversationIDs, err := l.svcCtx.ConversationDB.FindUserIDAllConversationID(l.ctx, req.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversation IDs: %w", err)
	}

	// 计算哈希：排序后使用 MD5，取前8个字符转换为 uint64
	hash := calculateConversationIDsHash(conversationIDs)

	resp.Hash = hash
	return resp, nil
}

// calculateConversationIDsHash 计算会话ID列表的哈希值
func calculateConversationIDsHash(conversationIDs []string) uint64 {
	if len(conversationIDs) == 0 {
		return 0
	}

	// 使用 JSON marshal（与 open-im-server 保持一致）
	data, _ := json.Marshal(conversationIDs)

	// MD5 哈希
	hash := md5.Sum(data)

	// 取前8个字节转换为 uint64
	return binary.BigEndian.Uint64(hash[:8])
}
