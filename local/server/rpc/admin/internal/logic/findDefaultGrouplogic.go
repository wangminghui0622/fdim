package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDefaultGroupLogic {
	return &FindDefaultGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindDefaultGroupLogic) FindDefaultGroup(req *admin.FindDefaultGroupReq) (*admin.FindDefaultGroupResp, error) {
	// FindDefaultGroupReq 没有 GroupIDs 字段，返回所有默认群组
	groupIDs, err := l.svcCtx.AdminDB.FindDefaultGroup(l.ctx, nil)
	if err != nil {
		l.Errorf("FindDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find default groups")
	}

	return &admin.FindDefaultGroupResp{
		GroupIDs: groupIDs,
	}, nil
}
