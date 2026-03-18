package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type DelUserClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelUserClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserClientConfigLogic {
	return &DelUserClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelUserClientConfigLogic) DelUserClientConfig(req *types.DelUserClientConfigReq) (resp *types.DelUserClientConfigResp, err error) {
	// 转换请求参数
	rpcReq := &user.DelUserClientConfigReq{
		UserID: req.UserID,
		Keys:   req.Keys,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.DelUserClientConfig(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DelUserClientConfigResp{
	}

	return resp, nil
}
