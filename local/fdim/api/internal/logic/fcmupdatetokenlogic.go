package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type FcmUpdateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFcmUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FcmUpdateTokenLogic {
	return &FcmUpdateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FcmUpdateTokenLogic) FcmUpdateToken(req *types.FcmUpdateTokenReq) (resp *types.FcmUpdateTokenResp, err error) {
	// 转换请求参数
	rpcReq := &third.FcmUpdateTokenReq{
		PlatformID: req.PlatformID,
		FcmToken:   req.FcmToken,
		Account:    req.Account,
		ExpireTime: req.ExpireTime,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ThirdClient.FcmUpdateToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.FcmUpdateTokenResp{
	}

	return resp, nil
}
