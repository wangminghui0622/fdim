package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/auth"
)

type ParseTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParseTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseTokenLogic {
	return &ParseTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseTokenLogic) ParseToken(req *types.ParseTokenReq) (resp *types.ParseTokenResp, err error) {
	// 转换请求参数
	rpcReq := &auth.ParseTokenReq{
		Token: req.Token,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.AuthClient.ParseToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ParseTokenResp{
		UserID:            rpcResp.UserID,
		PlatformID:        rpcResp.PlatformID,
		ExpireTimeSeconds: rpcResp.ExpireTimeSeconds,
	}

	return resp, nil
}
