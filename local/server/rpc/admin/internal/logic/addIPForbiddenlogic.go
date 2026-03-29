package logic

import (
	"context"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddIPForbiddenLogic {
	return &AddIPForbiddenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddIPForbiddenLogic) AddIPForbidden(req *admin.AddIPForbiddenReq) (*admin.AddIPForbiddenResp, error) {
	// 1. ֤
	if len(req.Forbiddens) == 0 {
		return nil, errs.ErrArgs.WrapMsg("forbiddens cannot be empty")
	}

	// 2. IPֹб
	forbiddens := make([]*database.IPForbidden, len(req.Forbiddens))
	now := time.Now()
	for i, fb := range req.Forbiddens {
		forbiddens[i] = &database.IPForbidden{
			IP:            fb.Ip,
			LimitRegister: fb.LimitRegister,
			LimitLogin:    fb.LimitLogin,
			CreateTime:    now,
		}
	}

	// 3. IPֹ
	if err := l.svcCtx.AdminDB.AddIPForbidden(l.ctx, forbiddens); err != nil {
		l.Errorf("AddIPForbidden failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add IP forbidden")
	}

	l.Infof("Added IP forbidden: count=%d", len(req.Forbiddens))
	return &admin.AddIPForbiddenResp{}, nil
}
