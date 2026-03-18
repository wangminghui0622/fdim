package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelInvitationCodeLogic {
	return &DelInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelInvitationCodeLogic) DelInvitationCode(req *admin.DelInvitationCodeReq) (*admin.DelInvitationCodeResp, error) {
	// 1. 验证参数
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes cannot be empty")
	}

	// 2. 删除邀请码
	if err := l.svcCtx.AdminDB.DelInvitationCode(l.ctx, req.Codes); err != nil {
		l.Errorf("DelInvitationCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete invitation codes")
	}

	l.Infof("Deleted invitation codes: count=%d", len(req.Codes))
	return &admin.DelInvitationCodeResp{}, nil
}
