package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMaxSeqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMaxSeqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMaxSeqsLogic {
	return &GetMaxSeqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMaxSeqsLogic) GetMaxSeqs(req *msg.GetMaxSeqsReq) (*msg.SeqsInfoResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.SeqsInfoResp{
		MaxSeqs: make(map[string]int64),
	}

	for _, convID := range req.ConversationIDs {
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		resp.MaxSeqs[convID] = maxSeq
	}

	return resp, nil
}
