package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDefaultGroupLogic {
	return &SearchDefaultGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchDefaultGroupLogic) SearchDefaultGroup(req *admin.SearchDefaultGroupReq) (*admin.SearchDefaultGroupResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索默认群组
	total, groupIDs, err := l.svcCtx.AdminDB.SearchDefaultGroup(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search default groups")
	}

	return &admin.SearchDefaultGroupResp{
		Total:    uint32(total),
		GroupIDs: groupIDs,
	}, nil
}
