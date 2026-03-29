package logic

import (
	"context"
	"strings"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsHasReadAndMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsHasReadAndMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsHasReadAndMaxSeqLogic {
	return &GetConversationsHasReadAndMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationsHasReadAndMaxSeqLogic) GetConversationsHasReadAndMaxSeq(req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID is required")
	}
	if len(req.ConversationIDs) == 0 {
		return &msg.GetConversationsHasReadAndMaxSeqResp{}, nil
	}
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.GetConversationsHasReadAndMaxSeqResp{
		Seqs: make(map[string]*msg.Seqs),
	}

	for _, convID := range req.ConversationIDs {
		if convID == "" {
			continue
		}
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		hasReadSeq, err := l.svcCtx.MsgCache.GetHasReadSeq(l.ctx, convID, req.UserID)
		if err != nil {
			l.Errorw("GetHasReadSeq failed", logx.Field("conversationID", convID), logx.Field("userID", req.UserID), logx.Field("error", err))
			continue
		}

		// 获取对方的已读状态（peerReadSeq）
		// 单聊会话格式: si_userA_userB，需要找到对方用户的 hasReadSeq
		var peerReadSeq int64 = 0
		if strings.HasPrefix(convID, "si_") {
			peerUserID := l.getPeerUserID(convID, req.UserID)
			if peerUserID != "" {
				peerReadSeq, _ = l.svcCtx.MsgCache.GetHasReadSeq(l.ctx, convID, peerUserID)
			}
		}

		resp.Seqs[convID] = &msg.Seqs{
			MaxSeq:      maxSeq,
			HasReadSeq:  hasReadSeq,
			PeerReadSeq: peerReadSeq,
		}
	}

	return resp, nil
}

// getPeerUserID 从单聊会话ID中提取对方用户ID
// 会话ID格式: si_userA_userB
func (l *GetConversationsHasReadAndMaxSeqLogic) getPeerUserID(convID, myUserID string) string {
	// 移除 "si_" 前缀
	if !strings.HasPrefix(convID, "si_") {
		return ""
	}
	suffix := strings.TrimPrefix(convID, "si_")
	// 格式: userA_userB
	parts := strings.Split(suffix, "_")
	if len(parts) != 2 {
		return ""
	}
	if parts[0] == myUserID {
		return parts[1]
	}
	if parts[1] == myUserID {
		return parts[0]
	}
	return ""
}
