package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ProcessUserCommandGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessUserCommandGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandGetLogic {
	return &ProcessUserCommandGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessUserCommandGetLogic) ProcessUserCommandGet(req *types.ProcessUserCommandGetReq) (resp *types.ProcessUserCommandGetResp, err error) {
	// 转换请求参数
	rpcReq := &user.ProcessUserCommandGetReq{
		UserID: req.UserID,
		Type:   req.Type,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.ProcessUserCommandGet(l.ctx, rpcReq)
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

	resp = &types.ProcessUserCommandGetResp{
		Commands: commands,
	}

	return resp, nil
}
