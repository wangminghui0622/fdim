package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/wrapperspb"
)

type ProcessUserCommandAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessUserCommandAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandAddLogic {
	return &ProcessUserCommandAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessUserCommandAddLogic) ProcessUserCommandAdd(req *types.ProcessUserCommandAddReq) (resp *types.ProcessUserCommandAddResp, err error) {
	// 转换请求参数
	rpcReq := &user.ProcessUserCommandAddReq{
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
	_, err = l.svcCtx.UserClient.ProcessUserCommandAdd(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ProcessUserCommandAddResp{
	}

	return resp, nil
}
