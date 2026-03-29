package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchBlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchBlockUserLogic {
	return &SearchBlockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchBlockUserLogic) SearchBlockUser(req *admin.SearchBlockUserReq) (*admin.SearchBlockUserResp, error) {
	// 1. ֤ҳ
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. û
	total, blocks, err := l.svcCtx.AdminDB.SearchBlockUser(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchBlockUser failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search block users")
	}

	// 3. תΪӦʽ
	results := make([]*admin.BlockUserInfo, 0, len(blocks))
	for _, block := range blocks {
		results = append(results, &admin.BlockUserInfo{
			UserID:     block.UserID,
			Reason:     block.Reason,
			OpUserID:   block.OperatorUserID,
			CreateTime: block.CreateTime.UnixMilli(),
		})
	}

	return &admin.SearchBlockUserResp{
		Total: uint32(total),
		Users: results,
	}, nil
}
