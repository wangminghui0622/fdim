package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinGroupLogic) JoinGroup(req *types.JoinGroupReq) (resp *types.JoinGroupResp, err error) {
	rpcReq := &user.JoinGroupReq{
		GroupID:       req.GroupID,
		ReqMessage:    req.ReqMessage,
		JoinSource:    req.JoinSource,
		InviterUserID: req.InviterUserID,
		Ex:            req.Ex,
	}

	_, err = l.svcCtx.GroupClient.JoinGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.JoinGroupResp{
	}

	return resp, nil
}
