package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchNotificationAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchNotificationAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchNotificationAccountLogic {
	return &SearchNotificationAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchNotificationAccountLogic) SearchNotificationAccount(req *user.SearchNotificationAccountReq) (*user.SearchNotificationAccountResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// TODO: 实现搜索通知账户的逻辑
	// 需要实?FindByNickname, FindNotification, FindSystemAccount 等方?
	return nil, fmt.Errorf("SearchNotificationAccount not fully implemented yet")
}
