package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsFriendLogic) IsFriend(req *types.IsFriendReq) (resp *types.IsFriendResp, err error) {
	rpcReq := &user.IsFriendReq{
		UserID1: req.UserID1,
		UserID2: req.UserID2,
	}

	rpcResp, err := l.svcCtx.FriendClient.IsFriend(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.IsFriendResp{
		InUser1Friends: rpcResp.InUser1Friends,
		InUser2Friends: rpcResp.InUser2Friends,
	}

	return resp, nil
}
