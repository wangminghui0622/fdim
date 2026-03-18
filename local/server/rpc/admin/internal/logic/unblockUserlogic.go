package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnblockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnblockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockUserLogic {
	return &UnblockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnblockUserLogic) UnblockUser(req *admin.UnblockUserReq) (*admin.UnblockUserResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs cannot be empty")
	}

	// 2. 删除封禁记录
	if err := l.svcCtx.AdminDB.DelBlockUser(l.ctx, req.UserIDs); err != nil {
		l.Errorf("DelBlockUser failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to unblock users")
	}

	l.Infof("Unblocked users: userIDs=%v", req.UserIDs)
	return &admin.UnblockUserResp{}, nil
}
