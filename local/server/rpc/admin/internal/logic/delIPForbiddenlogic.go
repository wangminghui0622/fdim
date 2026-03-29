package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelIPForbiddenLogic {
	return &DelIPForbiddenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelIPForbiddenLogic) DelIPForbidden(req *admin.DelIPForbiddenReq) (*admin.DelIPForbiddenResp, error) {
	// 1. ֤
	if len(req.Ips) == 0 {
		return nil, errs.ErrArgs.WrapMsg("ips cannot be empty")
	}

	// 2. ɾIPֹ
	if err := l.svcCtx.AdminDB.DelIPForbidden(l.ctx, req.Ips); err != nil {
		l.Errorf("DelIPForbidden failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete IP forbidden")
	}

	l.Infof("Deleted IP forbidden: ips=%v", req.Ips)
	return &admin.DelIPForbiddenResp{}, nil
}
