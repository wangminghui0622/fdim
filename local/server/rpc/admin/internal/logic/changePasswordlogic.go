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
	// 1.  context лȡ userIDͨм token н context
	userID := mcontext.GetOpUserID(l.ctx)
	if userID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID not found in context")
	}

	// 2. ֤벻Ϊ
	if req.Password == "" {
		return nil, errs.ErrArgs.WrapMsg("new password cannot be empty")
	}

	// 3. ֤ԱǷ
	adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("GetAdminUserID failed: %v", err)
		return nil, errs.ErrArgs.WrapMsg("admin not found")
	}

	// 4. 
	if err := l.svcCtx.AdminDB.ChangePassword(l.ctx, userID, req.Password); err != nil {
		l.Errorf("ChangePassword failed: %v", err)
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	l.Infof("Admin %s (userID: %s) changed password successfully", adminInfo.Account, userID)

	return &admin.ChangePasswordResp{}, nil
}
