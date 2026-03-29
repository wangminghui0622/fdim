package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
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

func (l *ResetPasswordLogic) ResetPassword(req *chat.ResetPasswordReq) (*chat.ResetPasswordResp, error) {
	// 1. ֤
	if req.Password == "" {
		return nil, errs.ErrArgs.WrapMsg("password cannot be empty")
	}
	if req.VerifyCode == "" {
		return nil, errs.ErrArgs.WrapMsg("verify code cannot be empty")
	}

	// 2. ֻŻûУ֤
	var userAccount *database.UserAccount
	var err error

	if req.PhoneNumber != "" {
		if req.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code is required for phone number")
		}
		if !strings.HasPrefix(req.AreaCode, "+") {
			req.AreaCode = "+" + req.AreaCode
		}
		if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		userAccount, err = l.svcCtx.ChatDB.GetUserAccountByPhone(l.ctx, req.AreaCode, req.PhoneNumber)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("user not found")
		}
		// ֤֤루ֻ
		verifyCode, err := l.svcCtx.ChatDB.FindVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode, req.VerifyCode)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid verify code")
		}
		// ֤Ƿ
		if time.Now().After(verifyCode.ExpireTime) {
			return nil, errs.ErrArgs.WrapMsg("verify code expired")
		}
		// ɾʹõ֤
		_ = l.svcCtx.ChatDB.DelVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode)
	} else if req.Email != "" {
		userAccount, err = l.svcCtx.ChatDB.GetUserAccountByEmail(l.ctx, req.Email)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("user not found")
		}
		// ֤֤루䣩
		verifyCode, err := l.svcCtx.ChatDB.FindVerifyCodeByEmail(l.ctx, req.Email, req.VerifyCode)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid verify code")
		}
		if time.Now().After(verifyCode.ExpireTime) {
			return nil, errs.ErrArgs.WrapMsg("verify code expired")
		}
		_ = l.svcCtx.ChatDB.DelVerifyCodeByEmail(l.ctx, req.Email)
	} else {
		return nil, errs.ErrArgs.WrapMsg("phone number or email must be set")
	}

	// 3. 
	if err := l.svcCtx.ChatDB.UpdateUserAccount(l.ctx, userAccount.UserID, map[string]interface{}{
		"password": req.Password,
	}); err != nil {
		l.Errorf("UpdateUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to reset password")
	}

	l.Infof("Password reset successfully: userID=%s", userAccount.UserID)
	return &chat.ResetPasswordResp{}, nil
}
