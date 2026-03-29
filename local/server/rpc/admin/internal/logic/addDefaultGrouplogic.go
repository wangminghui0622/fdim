package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDefaultGroupLogic {
	return &AddDefaultGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddDefaultGroupLogic) AddDefaultGroup(req *admin.AddDefaultGroupReq) (*admin.AddDefaultGroupResp, error) {
	// 1. ֤
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("groupIDs cannot be empty")
	}

	// 2. ǷѴ
	exists, err := l.svcCtx.AdminDB.FindDefaultGroup(l.ctx, req.GroupIDs)
	if err != nil {
		l.Errorf("FindDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check default groups")
	}
	if len(exists) > 0 {
		return nil, fmt.Errorf("some groupIDs already exist as default groups: %v", exists)
	}

	// 3. ĬȺ
	if err := l.svcCtx.AdminDB.AddDefaultGroup(l.ctx, req.GroupIDs); err != nil {
		l.Errorf("AddDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add default groups")
	}

	l.Infof("Added default groups: groupIDs=%v", req.GroupIDs)
	return &admin.AddDefaultGroupResp{}, nil
}
