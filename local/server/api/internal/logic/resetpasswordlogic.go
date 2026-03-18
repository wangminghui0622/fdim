package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (*types.ResetPasswordResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	_, err := l.svcCtx.ChatClient.ResetPassword(l.ctx, &chat.ResetPasswordReq{
		AreaCode:    req.AreaCode,
		PhoneNumber: req.PhoneNumber,
		VerifyCode:  req.VerifyCode,
		Password:    req.NewPassword,
		Email:       req.Email,
	})
	if err != nil {
		l.Errorf("ResetPassword failed: %v", err)
		return nil, err
	}

	return &types.ResetPasswordResp{}, nil
}
