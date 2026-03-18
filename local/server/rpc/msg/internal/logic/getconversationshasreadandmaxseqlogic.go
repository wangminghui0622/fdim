package logic

import (
	"context"

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

		resp.Seqs[convID] = &msg.Seqs{
			MaxSeq:     maxSeq,
			HasReadSeq: hasReadSeq,
			// MaxSeqTime 暂不从存储中反查，先置 0
		}
	}

	// pinnedConversationIDs 暂不实现，保持为空
	return resp, nil
}
