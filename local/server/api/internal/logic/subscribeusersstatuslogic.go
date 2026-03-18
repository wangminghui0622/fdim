package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type SubscribeUsersStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeUsersStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeUsersStatusLogic {
	return &SubscribeUsersStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeUsersStatusLogic) SubscribeUsersStatus(req *types.SubscribeOrCancelUsersStatusReq) (resp *types.SubscribeOrCancelUsersStatusResp, err error) {
	// 转换请求参数
	rpcReq := &user.SubscribeOrCancelUsersStatusReq{
		UserID:  req.UserID,
		UserIDs: req.UserIDs,
		Genre:   req.Genre,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.SubscribeOrCancelUsersStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SubscribeOrCancelUsersStatusResp{
	}

	return resp, nil
}
