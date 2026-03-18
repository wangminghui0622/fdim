package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type QuitGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitGroupLogic {
	return &QuitGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuitGroupLogic) QuitGroup(req *types.QuitGroupReq) (resp *types.QuitGroupResp, err error) {
	rpcReq := &user.QuitGroupReq{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}

	_, err = l.svcCtx.GroupClient.QuitGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.QuitGroupResp{
	}

	return resp, nil
}
