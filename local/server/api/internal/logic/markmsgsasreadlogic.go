package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type MarkMsgsAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkMsgsAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkMsgsAsReadLogic {
	return &MarkMsgsAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkMsgsAsReadLogic) MarkMsgsAsRead(req *types.MarkMsgsAsReadReq) (resp *types.MarkMsgsAsReadResp, err error) {
	// 转换请求参数
	rpcReq := &msg.MarkMsgsAsReadReq{
		ConversationID: req.ConversationID,
		Seqs:           req.Seqs,
		UserID:         req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.MarkMsgsAsRead(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.MarkMsgsAsReadResp{
	}

	return resp, nil
}
