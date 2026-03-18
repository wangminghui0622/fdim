package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDefaultGroupLogic {
	return &DelDefaultGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelDefaultGroupLogic) DelDefaultGroup(req *admin.DelDefaultGroupReq) (*admin.DelDefaultGroupResp, error) {
	// 1. 验证参数
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("groupIDs cannot be empty")
	}

	// 2. 检查是否存在
	exists, err := l.svcCtx.AdminDB.FindDefaultGroup(l.ctx, req.GroupIDs)
	if err != nil {
		l.Errorf("FindDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check default groups")
	}
	if len(exists) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("no default groups found for the given groupIDs")
	}

	// 3. 删除默认群组
	if err := l.svcCtx.AdminDB.DelDefaultGroup(l.ctx, req.GroupIDs); err != nil {
		l.Errorf("DelDefaultGroup failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete default groups")
	}

	l.Infof("Deleted default groups: groupIDs=%v", req.GroupIDs)
	return &admin.DelDefaultGroupResp{}, nil
}
