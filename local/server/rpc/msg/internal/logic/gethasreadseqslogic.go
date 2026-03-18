package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetHasReadSeqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHasReadSeqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHasReadSeqsLogic {
	return &GetHasReadSeqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetHasReadSeqsLogic) GetHasReadSeqs(req *msg.GetHasReadSeqsReq) (*msg.SeqsInfoResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.SeqsInfoResp{
		MaxSeqs: make(map[string]int64),
	}

	for _, convID := range req.ConversationIDs {
		hasReadSeq, err := l.svcCtx.MsgCache.GetHasReadSeq(l.ctx, convID, req.UserID)
		if err != nil {
			l.Errorw("GetHasReadSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		resp.MaxSeqs[convID] = hasReadSeq
	}

	return resp, nil
}
