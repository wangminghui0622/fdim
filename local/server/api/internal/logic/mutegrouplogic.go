package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type MuteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMuteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupLogic {
	return &MuteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MuteGroupLogic) MuteGroup(req *types.MuteGroupReq) (resp *types.MuteGroupResp, err error) {
	rpcReq := &user.MuteGroupReq{
		GroupID: req.GroupID,
	}

	_, err = l.svcCtx.GroupClient.MuteGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.MuteGroupResp{
	}

	return resp, nil
}
