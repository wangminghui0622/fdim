package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserIPLimitLoginLogic {
	return &SearchUserIPLimitLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserIPLimitLoginLogic) SearchUserIPLimitLogin(req *admin.SearchUserIPLimitLoginReq) (*admin.SearchUserIPLimitLoginResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索用户IP登录限制
	total, limits, err := l.svcCtx.AdminDB.SearchUserLimitLogin(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchUserLimitLogin failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search user IP limit login")
	}

	// 3. 转换为响应格式
	results := make([]*admin.LimitUserLoginIP, 0, len(limits))
	for _, limit := range limits {
		results = append(results, &admin.LimitUserLoginIP{
			UserID:     limit.UserID,
			Ip:         limit.IP,
			CreateTime: limit.CreateTime.UnixMilli(),
		})
	}

	return &admin.SearchUserIPLimitLoginResp{
		Total:  uint32(total),
		Limits: results,
	}, nil
}
