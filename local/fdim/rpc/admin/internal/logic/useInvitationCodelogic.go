package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UseInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseInvitationCodeLogic {
	return &UseInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UseInvitationCodeLogic) UseInvitationCode(req *admin.UseInvitationCodeReq) (*admin.UseInvitationCodeResp, error) {
	// 1. 验证参数
	if req.Code == "" {
		return nil, errs.ErrArgs.WrapMsg("code cannot be empty")
	}
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID cannot be empty")
	}

	// 2. 检查邀请码是否存在
	invitations, err := l.svcCtx.AdminDB.FindInvitationRegister(l.ctx, []string{req.Code})
	if err != nil {
		l.Errorf("FindInvitationRegister failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find invitation code")
	}
	if len(invitations) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("invitation code not found")
	}

	// 3. 检查邀请码是否已被使用
	inv := invitations[0]
	if inv.UsedByUserID != "" {
		return nil, fmt.Errorf("invitation code already used by user: %s", inv.UsedByUserID)
	}

	// 4. 标记邀请码为已使用
	if err := l.svcCtx.AdminDB.UseInvitationCode(l.ctx, req.Code, req.UserID); err != nil {
		l.Errorf("UseInvitationCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to use invitation code")
	}

	l.Infof("User %s used invitation code: %s", req.UserID, req.Code)
	return &admin.UseInvitationCodeResp{}, nil
}
