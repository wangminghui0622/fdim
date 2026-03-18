package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFriendLogic) IsFriend(req *user.IsFriendReq) (*user.IsFriendResp, error) {
	resp := &user.IsFriendResp{}

	// Ȩ����֤
	if err := authverify.CheckAccessIn(l.ctx, req.UserID1, req.UserID2); err != nil {
		return nil, err
	}

	inUser1Friends, inUser2Friends, err := l.svcCtx.FriendDB.CheckIn(l.ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}

	resp.InUser1Friends = inUser1Friends
	resp.InUser2Friends = inUser2Friends
	return resp, nil
}
