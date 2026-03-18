package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type DeleteMsgPhysicalBySeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMsgPhysicalBySeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgPhysicalBySeqLogic {
	return &DeleteMsgPhysicalBySeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMsgPhysicalBySeqLogic) DeleteMsgPhysicalBySeq(req *types.DeleteMsgPhysicalBySeqReq) (resp *types.DeleteMsgPhysicalBySeqResp, err error) {
	// 转换请求参数
	rpcReq := &msg.DeleteMsgPhysicalBySeqReq{
		ConversationID: req.ConversationID,
		Seqs:           []int64{req.Seq},
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.DeleteMsgPhysicalBySeq(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DeleteMsgPhysicalBySeqResp{
	}

	return resp, nil
}
