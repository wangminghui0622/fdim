package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type SetUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserClientConfigLogic {
	return &SetUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUserClientConfigLogic) SetUserClientConfig(req *types.SetUserClientConfigReq) (resp *types.SetUserClientConfigResp, err error) {
	// 转换请求参数
	rpcReq := &user.SetUserClientConfigReq{
		UserID:  req.UserID,
		Configs: req.Configs,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.SetUserClientConfig(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SetUserClientConfigResp{
	}

	return resp, nil
}
