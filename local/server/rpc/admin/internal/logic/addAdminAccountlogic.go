package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminAccountLogic {
	return &AddAdminAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAdminAccountLogic) AddAdminAccount(req *admin.AddAdminAccountReq) (*admin.AddAdminAccountResp, error) {
	// 1. ֤
	if req.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("account is empty")
	}
	if req.Password == "" {
		return nil, errs.ErrArgs.WrapMsg("password is empty")
	}

	// 2. ˻ǷѴ
	_, err := l.svcCtx.AdminDB.GetAdmin(l.ctx, req.Account)
	if err == nil {
		return nil, errs.ErrArgs.WrapMsg("account already exists")
	}

	// 3.  userIDʹ UUID
	userID := uuid.New().String()

	// 4. Ա
	adminAccount := &database.Admin{
		Account:    req.Account,
		Password:   req.Password,
		FaceURL:    req.FaceURL,
		Nickname:   req.Nickname,
		UserID:     userID,
		Level:      1, // Ĭϼ
		CreateTime: time.Now().Unix(),
	}

	// 5. ӵݿ
	if err := l.svcCtx.AdminDB.AddAdminAccount(l.ctx, []*database.Admin{adminAccount}); err != nil {
		l.Errorf("AddAdminAccount failed: %v", err)
		return nil, fmt.Errorf("failed to add admin account: %w", err)
	}

	l.Infof("Added admin account: account=%s, userID=%s", req.Account, userID)

	return &admin.AddAdminAccountResp{}, nil
}
