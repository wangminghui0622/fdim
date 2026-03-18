package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type GetSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSeqLogic {
	return &GetSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSeqLogic) GetSeq(req *types.GetSeqReq) (resp *types.GetSeqResp, err error) {
	if req.ConversationID == "" {
		return &types.GetSeqResp{MaxSeq: 0}, nil
	}

	// 使用已实现的 GetConversationsHasReadAndMaxSeq RPC
	rpcResp, err := l.svcCtx.MsgClient.GetConversationsHasReadAndMaxSeq(l.ctx, &msg.GetConversationsHasReadAndMaxSeqReq{
		UserID:          req.UserID,
		ConversationIDs: []string{req.ConversationID},
	})
	if err != nil {
		return nil, err
	}

	var maxSeq int64
	if seqs, ok := rpcResp.Seqs[req.ConversationID]; ok {
		maxSeq = seqs.MaxSeq
	}

	return &types.GetSeqResp{
		MaxSeq: maxSeq,
	}, nil
}
