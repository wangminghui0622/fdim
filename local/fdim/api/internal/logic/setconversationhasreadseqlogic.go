package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type SetConversationHasReadSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetConversationHasReadSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationHasReadSeqLogic {
	return &SetConversationHasReadSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetConversationHasReadSeqLogic) SetConversationHasReadSeq(req *types.SetConversationHasReadSeqReq) (resp *types.SetConversationHasReadSeqResp, err error) {
	// 转换请求参数
	rpcReq := &msg.SetConversationHasReadSeqReq{
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		HasReadSeq:     req.HasReadSeq,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.SetConversationHasReadSeq(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SetConversationHasReadSeqResp{
	}

	return resp, nil
}
