package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAdminAccountLogic {
	return &DelAdminAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelAdminAccountLogic) DelAdminAccount(req *admin.DelAdminAccountReq) (*admin.DelAdminAccountResp, error) {
	// 1. ֤
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs is empty")
	}

	// 2. Ƿظ userID
	userIDMap := make(map[string]bool)
	for _, userID := range req.UserIDs {
		if userIDMap[userID] {
			return nil, errs.ErrArgs.WrapMsg("user ids is duplicate")
		}
		userIDMap[userID] = true
	}

	// 3. ǷǳԱɾԱ
	// Ա 100 (AdvancedUserLevel)
	const AdvancedUserLevel = 100
	for _, userID := range req.UserIDs {
		adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, userID)
		if err != nil {
			l.Errorf("GetAdminUserID failed: %v", err)
			return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("admin %s not found", userID))
		}
		if adminInfo.Level >= AdvancedUserLevel {
			return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("%s is superAdminID, cannot delete", userID))
		}
	}

	// 4. ɾԱ˻
	if err := l.svcCtx.AdminDB.DelAdminAccount(l.ctx, req.UserIDs); err != nil {
		l.Errorf("DelAdminAccount failed: %v", err)
		return nil, fmt.Errorf("failed to delete admin account: %w", err)
	}

	l.Infof("Deleted admin accounts: userIDs=%v", req.UserIDs)

	return &admin.DelAdminAccountResp{}, nil
}
