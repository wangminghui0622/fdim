package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type GetConversationsHasReadAndMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationsHasReadAndMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsHasReadAndMaxSeqLogic {
	return &GetConversationsHasReadAndMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsHasReadAndMaxSeqLogic) GetConversationsHasReadAndMaxSeq(req *types.GetConversationsHasReadAndMaxSeqReq) (resp *types.GetConversationsHasReadAndMaxSeqResp, err error) {
	// 转换请求参数
	rpcReq := &msg.GetConversationsHasReadAndMaxSeqReq{
		UserID:          req.UserID,
		ConversationIDs: req.ConversationIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.GetConversationsHasReadAndMaxSeq(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	seqs := make(map[string]interface{})
	for k, v := range rpcResp.Seqs {
		seqs[k] = v
	}

	resp = &types.GetConversationsHasReadAndMaxSeqResp{
		Seqs: seqs,
	}

	return resp, nil
}
