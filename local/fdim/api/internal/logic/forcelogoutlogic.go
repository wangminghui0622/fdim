package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/auth"
)

type ForceLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForceLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForceLogoutLogic {
	return &ForceLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForceLogoutLogic) ForceLogout(req *types.ForceLogoutReq) (resp *types.ForceLogoutResp, err error) {
	// 转换请求参数
	rpcReq := &auth.ForceLogoutReq{
		PlatformID: req.PlatformID,
		UserID:     req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.AuthClient.ForceLogout(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ForceLogoutResp{
	}

	return resp, nil
}
