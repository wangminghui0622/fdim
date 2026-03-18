package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ProcessUserCommandDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessUserCommandDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandDeleteLogic {
	return &ProcessUserCommandDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessUserCommandDeleteLogic) ProcessUserCommandDelete(req *types.ProcessUserCommandDeleteReq) (resp *types.ProcessUserCommandDeleteResp, err error) {
	// 转换请求参数
	rpcReq := &user.ProcessUserCommandDeleteReq{
		UserID: req.UserID,
		Type:   req.Type,
		Uuid:   req.Uuid,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.ProcessUserCommandDelete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ProcessUserCommandDeleteResp{
	}

	return resp, nil
}
