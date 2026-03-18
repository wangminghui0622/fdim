package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/auth"
)

type GetAdminTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminTokenLogic {
	return &GetAdminTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminTokenLogic) GetAdminToken(req *types.GetAdminTokenReq) (resp *types.GetAdminTokenResp, err error) {
	// 转换请求参数
	rpcReq := &auth.GetAdminTokenReq{
		Secret: req.Secret,
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.AuthClient.GetAdminToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetAdminTokenResp{
		Token:             rpcResp.Token,
		ExpireTimeSeconds: rpcResp.ExpireTimeSeconds,
	}

	return resp, nil
}
