package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *chat.ChangePasswordReq) (*chat.ChangePasswordResp, error) {
	// 1. ֤
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user ID cannot be empty")
	}
	if req.CurrentPassword == "" {
		return nil, errs.ErrArgs.WrapMsg("current password cannot be empty")
	}
	if req.NewPassword == "" {
		return nil, errs.ErrArgs.WrapMsg("new password cannot be empty")
	}

	// 2. ȡû˻
	userAccount, err := l.svcCtx.ChatDB.GetUserAccountByUserID(l.ctx, req.UserID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("user not found")
	}

	// 3. ֤ǰ
	if userAccount.Password != req.CurrentPassword {
		return nil, errs.ErrArgs.WrapMsg("current password is incorrect")
	}

	// 4. 
	if err := l.svcCtx.ChatDB.UpdateUserAccount(l.ctx, req.UserID, map[string]interface{}{
		"password": req.NewPassword,
	}); err != nil {
		l.Errorf("UpdateUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to update password")
	}

	l.Infof("Password changed successfully: userID=%s", req.UserID)
	return &chat.ChangePasswordResp{}, nil
}
