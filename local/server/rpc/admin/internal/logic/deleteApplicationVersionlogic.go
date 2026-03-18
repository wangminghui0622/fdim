package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationVersionLogic {
	return &DeleteApplicationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteApplicationVersionLogic) DeleteApplicationVersion(req *admin.DeleteApplicationVersionReq) (*admin.DeleteApplicationVersionResp, error) {
	// 1. 验证参数
	if len(req.Id) == 0 {
		return nil, errs.ErrArgs.WrapMsg("id cannot be empty")
	}

	// 2. 删除应用版本
	if err := l.svcCtx.AdminDB.DeleteApplicationVersion(l.ctx, req.Id); err != nil {
		l.Errorf("DeleteApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete application versions")
	}

	l.Infof("Deleted application versions: count=%d", len(req.Id))
	return &admin.DeleteApplicationVersionResp{}, nil
}
