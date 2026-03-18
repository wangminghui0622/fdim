package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type GetServerTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerTimeLogic {
	return &GetServerTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerTimeLogic) GetServerTime(req *types.GetServerTimeReq) (resp *types.GetServerTimeResp, err error) {
	// 转换请求参数
	rpcReq := &msg.GetServerTimeReq{}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.GetServerTime(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetServerTimeResp{
		ServerTime: rpcResp.ServerTime,
	}

	return resp, nil
}
