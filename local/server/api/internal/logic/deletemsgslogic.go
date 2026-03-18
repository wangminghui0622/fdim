package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type DeleteMsgsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgsLogic {
	return &DeleteMsgsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMsgsLogic) DeleteMsgs(req *types.DeleteMsgsReq) (resp *types.DeleteMsgsResp, err error) {
	// 转换请求参数
	rpcReq := &msg.DeleteMsgsReq{
		ConversationID: req.ConversationID,
		Seqs:           req.Seqs,
		UserID:         req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.DeleteMsgs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DeleteMsgsResp{
	}

	return resp, nil
}
