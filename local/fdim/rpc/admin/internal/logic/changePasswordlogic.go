package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/pkg/mcontext"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
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

func (l *ChangePasswordLogic) ChangePassword(req *admin.ChangePasswordReq) (*admin.ChangePasswordResp, error) {
	// 1. 从 context 中获取 userID（通常由中间件从 token 中解析并放入 context）
	userID := mcontext.GetOpUserID(l.ctx)
	if userID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID not found in context")
	}

	// 2. 验证新密码不能为空
	if req.Password == "" {
		return nil, errs.ErrArgs.WrapMsg("new password cannot be empty")
	}

	// 3. 验证管理员是否存在
	adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("GetAdminUserID failed: %v", err)
		return nil, errs.ErrArgs.WrapMsg("admin not found")
	}

	// 4. 更新密码
	if err := l.svcCtx.AdminDB.ChangePassword(l.ctx, userID, req.Password); err != nil {
		l.Errorf("ChangePassword failed: %v", err)
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	l.Infof("Admin %s (userID: %s) changed password successfully", adminInfo.Account, userID)

	return &admin.ChangePasswordResp{}, nil
}
