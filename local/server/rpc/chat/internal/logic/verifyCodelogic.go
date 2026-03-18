package logic

import (
	"context"
	"time"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
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

func (l *VerifyCodeLogic) VerifyCode(req *chat.VerifyCodeReq) (*chat.VerifyCodeResp, error) {
	// 1. 验证参数
	if req.VerifyCode == "" {
		return nil, errs.ErrArgs.WrapMsg("verify code cannot be empty")
	}

	// 2. 根据手机号或邮箱验证验证码
	if req.PhoneNumber != "" {
		if req.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code is required for phone number")
		}
		verifyCode, err := l.svcCtx.ChatDB.FindVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode, req.VerifyCode)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid verify code")
		}
		// 检查验证码是否过期
		if time.Now().After(verifyCode.ExpireTime) {
			return nil, errs.ErrArgs.WrapMsg("verify code expired")
		}
		// 验证成功后删除验证码
		_ = l.svcCtx.ChatDB.DelVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode)
	} else if req.Email != "" {
		verifyCode, err := l.svcCtx.ChatDB.FindVerifyCodeByEmail(l.ctx, req.Email, req.VerifyCode)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid verify code")
		}
		// 检查验证码是否过期
		if time.Now().After(verifyCode.ExpireTime) {
			return nil, errs.ErrArgs.WrapMsg("verify code expired")
		}
		// 验证成功后删除验证码
		_ = l.svcCtx.ChatDB.DelVerifyCodeByEmail(l.ctx, req.Email)
	} else {
		return nil, errs.ErrArgs.WrapMsg("phone number or email must be set")
	}

	l.Infof("Verify code verified successfully")
	return &chat.VerifyCodeResp{}, nil
}
