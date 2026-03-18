package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type UserClearAllMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClearAllMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClearAllMsgLogic {
	return &UserClearAllMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClearAllMsgLogic) UserClearAllMsg(req *types.UserClearAllMsgReq) (resp *types.UserClearAllMsgResp, err error) {
	// 转换请求参数
	rpcReq := &msg.UserClearAllMsgReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.UserClearAllMsg(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UserClearAllMsgResp{
	}

	return resp, nil
}
