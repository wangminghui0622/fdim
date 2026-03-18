package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsBlackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsBlackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsBlackLogic {
	return &IsBlackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsBlackLogic) IsBlack(req *user.IsBlackReq) (*user.IsBlackResp, error) {
	resp := &user.IsBlackResp{}

	// Ȩ����֤
	if err := authverify.CheckAccessIn(l.ctx, req.UserID1, req.UserID2); err != nil {
		return nil, err
	}

	// ��� UserID1 �Ƿ������� UserID2
	isBlack1, err := l.svcCtx.BlackDB.IsBlack(l.ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp.InUser1Blacks = isBlack1

	// ��� UserID2 �Ƿ������� UserID1
	isBlack2, err := l.svcCtx.BlackDB.IsBlack(l.ctx, req.UserID2, req.UserID1)
	if err != nil {
		return nil, err
	}
	resp.InUser2Blacks = isBlack2

	return resp, nil
}
