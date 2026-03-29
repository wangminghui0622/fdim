package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAppletLogic {
	return &DelAppletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelAppletLogic) DelApplet(req *admin.DelAppletReq) (*admin.DelAppletResp, error) {
	// 1. ֤
	if len(req.AppletIds) == 0 {
		return nil, errs.ErrArgs.WrapMsg("appletIds cannot be empty")
	}

	// 2. ɾС
	if err := l.svcCtx.AdminDB.DelApplet(l.ctx, req.AppletIds); err != nil {
		l.Errorf("DelApplet failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete applets")
	}

	l.Infof("Deleted applets: count=%d", len(req.AppletIds))
	return &admin.DelAppletResp{}, nil
}
