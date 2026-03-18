package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ProcessUserCommandGetAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessUserCommandGetAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandGetAllLogic {
	return &ProcessUserCommandGetAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessUserCommandGetAllLogic) ProcessUserCommandGetAll(req *types.ProcessUserCommandGetAllReq) (resp *types.ProcessUserCommandGetAllResp, err error) {
	// 转换请求参数
	rpcReq := &user.ProcessUserCommandGetAllReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.ProcessUserCommandGetAll(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var commands []types.UserCommand
	for _, cmd := range rpcResp.CommandResp {
		commands = append(commands, types.UserCommand{
			Uuid:  cmd.Uuid,
			Value: cmd.Value,
			Ex:    cmd.Ex,
		})
	}

	resp = &types.ProcessUserCommandGetAllResp{
		Commands: commands,
	}

	return resp, nil
}
