package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAdminAccountLogic {
	return &SearchAdminAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchAdminAccountLogic) SearchAdminAccount(req *admin.SearchAdminAccountReq) (*admin.SearchAdminAccountResp, error) {
	// 1. 验证参数
	if req.Pagination == nil {
		return nil, errs.ErrArgs.WrapMsg("pagination is empty")
	}
	if req.Pagination.ShowNumber == 0 {
		return nil, errs.ErrArgs.WrapMsg("showNumber is empty")
	}
	if req.Pagination.PageNumber == 0 {
		return nil, errs.ErrArgs.WrapMsg("pageNumber is empty")
	}

	// 2. 搜索管理员账户
	total, adminAccounts, err := l.svcCtx.AdminDB.SearchAdminAccount(l.ctx, req.Pagination)
	if err != nil {
		l.Errorf("SearchAdminAccount failed: %v", err)
		return nil, fmt.Errorf("failed to search admin account: %w", err)
	}

	// 3. 转换为响应格式
	accounts := make([]*admin.GetAdminInfoResp, 0, len(adminAccounts))
	for _, adminAccount := range adminAccounts {
		accounts = append(accounts, &admin.GetAdminInfoResp{
			Account:    adminAccount.Account,
			FaceURL:    adminAccount.FaceURL,
			Nickname:   adminAccount.Nickname,
			UserID:     adminAccount.UserID,
			Level:      adminAccount.Level,
			CreateTime: adminAccount.CreateTime,
		})
	}

	l.Infof("Searched admin accounts: total=%d, count=%d", total, len(accounts))

	return &admin.SearchAdminAccountResp{
		Total:         uint32(total),
		AdminAccounts: accounts,
	}, nil
}
