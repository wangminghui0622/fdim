package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type SetGlobalRecvMessageOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGlobalRecvMessageOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGlobalRecvMessageOptLogic {
	return &SetGlobalRecvMessageOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGlobalRecvMessageOptLogic) SetGlobalRecvMessageOpt(req *types.SetGlobalRecvMessageOptReq) (resp *types.SetGlobalRecvMessageOptResp, err error) {
	// 转换请求参数
	rpcReq := &user.SetGlobalRecvMessageOptReq{
		UserID:           req.UserID,
		GlobalRecvMsgOpt: req.GlobalRecvMsgOpt,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.SetGlobalRecvMessageOpt(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.SetGlobalRecvMessageOptResp{
	}

	return resp, nil
}
