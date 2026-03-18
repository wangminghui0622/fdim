package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindInvitationCodeLogic {
	return &FindInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindInvitationCodeLogic) FindInvitationCode(req *admin.FindInvitationCodeReq) (*admin.FindInvitationCodeResp, error) {
	// 1. 验证参数
	if len(req.Codes) == 0 {
		return nil, errs.ErrArgs.WrapMsg("codes cannot be empty")
	}

	// 2. 查找邀请码
	invitations, err := l.svcCtx.AdminDB.FindInvitationRegister(l.ctx, req.Codes)
	if err != nil {
		l.Errorf("FindInvitationRegister failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find invitation codes")
	}

	// 3. 转换为响应格式
	results := make([]*admin.InvitationRegister, 0, len(invitations))
	for _, inv := range invitations {
		results = append(results, &admin.InvitationRegister{
			InvitationCode: inv.InvitationCode,
			CreateTime:     inv.CreateTime.UnixMilli(),
			UsedUserID:     inv.UsedByUserID,
		})
	}

	return &admin.FindInvitationCodeResp{
		Codes: results,
	}, nil
}
