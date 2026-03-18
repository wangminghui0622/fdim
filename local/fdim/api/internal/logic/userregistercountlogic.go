package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type UserRegisterCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterCountLogic {
	return &UserRegisterCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterCountLogic) UserRegisterCount(req *types.UserRegisterCountReq) (resp *types.UserRegisterCountResp, err error) {
	// 转换请求参数
	rpcReq := &user.UserRegisterCountReq{
		Start: req.Start,
		End:   req.End,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.UserRegisterCount(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UserRegisterCountResp{
		Total:  rpcResp.Total,
		Before: rpcResp.Before,
	}

	return resp, nil
}
