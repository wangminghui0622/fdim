package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyCodeLogic) VerifyCode(req *types.VerifyCodeReq) (*types.VerifyCodeResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	_, err := l.svcCtx.ChatClient.VerifyCode(l.ctx, &chat.VerifyCodeReq{
		PhoneNumber: req.PhoneNumber,
		AreaCode:    req.AreaCode,
		Email:       req.Email,
		VerifyCode:  req.VerifyCode,
	})
	if err != nil {
		l.Errorf("VerifyCode failed: %v", err)
		return nil, err
	}

	return &types.VerifyCodeResp{}, nil
}
