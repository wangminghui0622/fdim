package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type RevokeMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeMsgLogic {
	return &RevokeMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeMsgLogic) RevokeMsg(req *types.RevokeMsgReq) (resp *types.RevokeMsgResp, err error) {
	// 转换请求参数
	rpcReq := &msg.RevokeMsgReq{
		ConversationID: req.ConversationID,
		Seq:            req.Seq,
		UserID:         req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.RevokeMsg(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.RevokeMsgResp{
	}

	return resp, nil
}
