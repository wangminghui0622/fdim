package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type CheckMsgIsSendSuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckMsgIsSendSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckMsgIsSendSuccessLogic {
	return &CheckMsgIsSendSuccessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckMsgIsSendSuccessLogic) CheckMsgIsSendSuccess(req *types.CheckMsgIsSendSuccessReq) (resp *types.CheckMsgIsSendSuccessResp, err error) {
	// 转换请求参数
	rpcReq := &msg.GetSendMsgStatusReq{}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.GetSendMsgStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.CheckMsgIsSendSuccessResp{
		Status: rpcResp.Status,
	}

	return resp, nil
}
