package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/wrapperspb"
)

type ProcessUserCommandUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessUserCommandUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandUpdateLogic {
	return &ProcessUserCommandUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessUserCommandUpdateLogic) ProcessUserCommandUpdate(req *types.ProcessUserCommandUpdateReq) (resp *types.ProcessUserCommandUpdateResp, err error) {
	// 转换请求参数
	rpcReq := &user.ProcessUserCommandUpdateReq{
		UserID: req.UserID,
		Type:   req.Type,
		Uuid:   req.Uuid,
	}

	if req.Value != nil {
		rpcReq.Value = wrapperspb.String(*req.Value)
	}
	if req.Ex != nil {
		rpcReq.Ex = wrapperspb.String(*req.Ex)
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.ProcessUserCommandUpdate(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ProcessUserCommandUpdateResp{
	}

	return resp, nil
}
