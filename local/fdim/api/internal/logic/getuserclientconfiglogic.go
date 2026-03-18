package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserClientConfigLogic {
	return &GetUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserClientConfigLogic) GetUserClientConfig(req *types.GetUserClientConfigReq) (resp *types.GetUserClientConfigResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetUserClientConfigReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetUserClientConfig(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetUserClientConfigResp{
		Configs: rpcResp.Configs,
	}

	return resp, nil
}
