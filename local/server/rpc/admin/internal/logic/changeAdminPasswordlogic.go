package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAdminPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAdminPasswordLogic {
	return &ChangeAdminPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeAdminPasswordLogic) ChangeAdminPassword(req *admin.ChangeAdminPasswordReq) (*admin.ChangeAdminPasswordResp, error) {
	// 1. ֤
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID is empty")
	}
	if req.CurrentPassword == "" {
		return nil, errs.ErrArgs.WrapMsg("currentPassword is empty")
	}
	if req.NewPassword == "" {
		return nil, errs.ErrArgs.WrapMsg("newPassword is empty")
	}
	if req.CurrentPassword == req.NewPassword {
		return nil, errs.ErrArgs.WrapMsg("currentPassword is equal to newPassword")
	}

	// 2. ȡԱϢ
	adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, req.UserID)
	if err != nil {
		l.Errorf("GetAdminUserID failed: %v", err)
		return nil, errs.ErrArgs.WrapMsg("admin not found")
	}

	// 3. ֤ǰ
	if adminInfo.Password != req.CurrentPassword {
		return nil, errs.ErrArgs.WrapMsg("password error")
	}

	// 4. 
	if err := l.svcCtx.AdminDB.ChangePassword(l.ctx, req.UserID, req.NewPassword); err != nil {
		l.Errorf("ChangePassword failed: %v", err)
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	l.Infof("Admin %s changed password successfully", req.UserID)

	return &admin.ChangeAdminPasswordResp{}, nil
}
