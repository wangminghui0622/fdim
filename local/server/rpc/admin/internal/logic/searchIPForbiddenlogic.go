package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchIPForbiddenLogic {
	return &SearchIPForbiddenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchIPForbiddenLogic) SearchIPForbidden(req *admin.SearchIPForbiddenReq) (*admin.SearchIPForbiddenResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索IP禁止
	total, forbiddens, err := l.svcCtx.AdminDB.SearchIPForbidden(l.ctx, req.Keyword, req.Status, req.Pagination)
	if err != nil {
		l.Errorf("SearchIPForbidden failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search IP forbidden")
	}

	// 3. 转换为响应格式
	results := make([]*admin.IPForbidden, 0, len(forbiddens))
	for _, fb := range forbiddens {
		results = append(results, &admin.IPForbidden{
			Ip:            fb.IP,
			LimitRegister: fb.LimitRegister,
			LimitLogin:    fb.LimitLogin,
			CreateTime:    fb.CreateTime.UnixMilli(),
		})
	}

	return &admin.SearchIPForbiddenResp{
		Total:      uint32(total),
		Forbiddens: results,
	}, nil
}
