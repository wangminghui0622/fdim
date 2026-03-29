package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddInvitationCodeLogic {
	return &AddInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddInvitationCodeLogic) AddInvitationCode(req *admin.AddInvitationCodeReq) (*admin.AddInvitationCodeResp, error) {
	// 1. ֤
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes cannot be empty")
	}

	// 2. ǷѴ
	exists, err := l.svcCtx.AdminDB.FindInvitationRegister(l.ctx, req.Codes)
	if err != nil {
		l.Errorf("FindInvitationRegister failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check invitation codes")
	}
	if len(exists) > 0 {
		existingCodes := make([]string, len(exists))
		for i, e := range exists {
			existingCodes[i] = e.InvitationCode
		}
		return nil, fmt.Errorf("some codes already exist: %v", existingCodes)
	}

	// 3. б
	invitations := make([]*database.InvitationRegister, len(req.Codes))
	now := time.Now()
	for i, code := range req.Codes {
		invitations[i] = &database.InvitationRegister{
			InvitationCode: code,
			UsedByUserID:   "",
			CreateTime:     now,
		}
	}

	// 4. 
	if err := l.svcCtx.AdminDB.AddInvitationCode(l.ctx, invitations); err != nil {
		l.Errorf("AddInvitationCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add invitation codes")
	}

	l.Infof("Added invitation codes: count=%d", len(req.Codes))
	return &admin.AddInvitationCodeResp{}, nil
}
