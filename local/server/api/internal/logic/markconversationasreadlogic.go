package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type MarkConversationAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkConversationAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkConversationAsReadLogic {
	return &MarkConversationAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkConversationAsReadLogic) MarkConversationAsRead(req *types.MarkConversationAsReadReq) (resp *types.MarkConversationAsReadResp, err error) {
	// 转换请求参数
	rpcReq := &msg.MarkConversationAsReadReq{
		ConversationID: req.ConversationID,
		UserID:         req.UserID,
		HasReadSeq:     req.HasReadSeq,
		Seqs:           req.Seqs,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.MarkConversationAsRead(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.MarkConversationAsReadResp{
	}

	return resp, nil
}
