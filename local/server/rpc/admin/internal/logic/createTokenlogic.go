package logic

import (
	"context"
	"fdim/pkg/errs"

	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTokenLogic {
	return &CreateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTokenLogic) CreateToken(req *admin.CreateTokenReq) (*admin.CreateTokenResp, error) {
	// 1.  Token
	tokenStr, expireDuration, err := l.svcCtx.Token.CreateToken(req.UserID, req.UserType)
	if err != nil {
		return nil, errs.WrapMsg(err, "CreateToken failed")
	}

	// 2.  Token  Redis
	expireSeconds := expireDuration.Milliseconds()
	if err := l.svcCtx.AdminDB.CacheToken(l.ctx, req.UserID, tokenStr, expireSeconds); err != nil {
		return nil, errs.WrapMsg(err, "CacheToken failed")
	}

	return &admin.CreateTokenResp{
		Token: tokenStr,
	}, nil
}
