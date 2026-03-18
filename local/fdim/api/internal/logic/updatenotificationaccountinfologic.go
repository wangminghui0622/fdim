package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type UpdateNotificationAccountInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateNotificationAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationAccountInfoLogic {
	return &UpdateNotificationAccountInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotificationAccountInfoLogic) UpdateNotificationAccountInfo(req *types.UpdateNotificationAccountInfoReq) (resp *types.UpdateNotificationAccountInfoResp, err error) {
	// 转换请求参数
	rpcReq := &user.UpdateNotificationAccountInfoReq{
		UserID:   req.UserID,
		FaceURL:  req.FaceURL,
		NickName: req.NickName,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.UpdateNotificationAccountInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UpdateNotificationAccountInfoResp{
	}

	return resp, nil
}
