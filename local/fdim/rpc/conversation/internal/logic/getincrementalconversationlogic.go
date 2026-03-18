package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// VersionSortChangeID is a constant for version tracking
const VersionSortChangeID = "version_sort_change"

type GetIncrementalConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalConversationLogic {
	return &GetIncrementalConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalConversationLogic) GetIncrementalConversation(req *conversation.GetIncrementalConversationReq) (*conversation.GetIncrementalConversationResp, error) {
	resp := &conversation.GetIncrementalConversationResp{}

	// 目前 go-im 的 ConversationDatabase 尚未实现版本日志表，为保持简单并保证正确性：
	// - 当客户端未携带 version（或 version=0）时，返回全量会话并标记 Full=true；
	// - 当携带 version>0 时，暂时仍返回全量（兼容行为），后续可接入完整 VersionLog 机制。

	conversations, err := l.svcCtx.ConversationDB.FindUserIDAllConversations(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations: %w", err)
	}

	resp.Insert = make([]*conversation.Conversation, 0, len(conversations))
	for _, conv := range conversations {
		resp.Insert = append(resp.Insert, &conversation.Conversation{
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

	// 简化版版本信息：使用一个递增的 uint 版本号（这里用会话数量近似），VersionID 使用固定标识
	resp.Version = uint64(len(conversations))
	resp.VersionID = VersionSortChangeID
	resp.Full = true
	resp.Delete = nil
	resp.Update = nil

	return resp, nil
}
