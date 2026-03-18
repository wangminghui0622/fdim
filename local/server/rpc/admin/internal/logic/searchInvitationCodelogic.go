package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchInvitationCodeLogic {
	return &SearchInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchInvitationCodeLogic) SearchInvitationCode(req *admin.SearchInvitationCodeReq) (*admin.SearchInvitationCodeResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索邀请码
	total, invitations, err := l.svcCtx.AdminDB.SearchInvitationCode(
		l.ctx,
		req.Status,
		req.UserIDs,
		req.Codes,
		req.Keyword,
		req.Pagination,
	)
	if err != nil {
		l.Errorf("SearchInvitationCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search invitation codes")
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

	return &admin.SearchInvitationCodeResp{
		Total: uint32(total),
		List:  results,
	}, nil
}
