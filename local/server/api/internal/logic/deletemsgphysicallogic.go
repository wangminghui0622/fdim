package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type DeleteMsgPhysicalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMsgPhysicalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgPhysicalLogic {
	return &DeleteMsgPhysicalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMsgPhysicalLogic) DeleteMsgPhysical(req *types.DeleteMsgPhysicalReq) (resp *types.DeleteMsgPhysicalResp, err error) {
	// 转换请求参数
	rpcReq := &msg.DeleteMsgPhysicalReq{
		ConversationIDs: []string{req.ConversationID},
		Timestamp:       0, // 如果需要时间戳，可以从请求中获?
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.DeleteMsgPhysical(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DeleteMsgPhysicalResp{
	}

	return resp, nil
}
