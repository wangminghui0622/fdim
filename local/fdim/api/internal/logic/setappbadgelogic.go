package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type SetAppBadgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetAppBadgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAppBadgeLogic {
	return &SetAppBadgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetAppBadgeLogic) SetAppBadge(req *types.SetAppBadgeReq) (resp *types.SetAppBadgeResp, err error) {
	// 转换请求参数
	rpcReq := &third.SetAppBadgeReq{
		UserID:         req.UserID,
		AppUnreadCount: req.AppUnreadCount,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ThirdClient.SetAppBadge(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SetAppBadgeResp{
	}

	return resp, nil
}
