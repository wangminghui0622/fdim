package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type AddNotificationAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddNotificationAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddNotificationAccountLogic {
	return &AddNotificationAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddNotificationAccountLogic) AddNotificationAccount(req *types.AddNotificationAccountReq) (resp *types.AddNotificationAccountResp, err error) {
	// 转换请求参数
	rpcReq := &user.AddNotificationAccountReq{
		UserID:         req.UserID,
		NickName:       req.NickName,
		FaceURL:        req.FaceURL,
		AppMangerLevel: req.AppMangerLevel,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.AddNotificationAccount(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.AddNotificationAccountResp{
		UserID:         rpcResp.UserID,
		FaceURL:        rpcResp.FaceURL,
		NickName:       rpcResp.NickName,
		AppMangerLevel: rpcResp.AppMangerLevel,
	}

	return resp, nil
}
