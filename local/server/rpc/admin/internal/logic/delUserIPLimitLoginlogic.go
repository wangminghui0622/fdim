package logic

import (
	"context"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserIPLimitLoginLogic {
	return &DelUserIPLimitLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserIPLimitLoginLogic) DelUserIPLimitLogin(req *admin.DelUserIPLimitLoginReq) (*admin.DelUserIPLimitLoginResp, error) {
	// 1. ֤
	if len(req.Limits) == 0 {
		return nil, errs.ErrArgs.WrapMsg("limits cannot be empty")
	}

	// 2. б
	limits := make([]*database.LimitUserLoginIP, len(req.Limits))
	for i, limit := range req.Limits {
		if limit.UserID == "" || limit.Ip == "" {
			return nil, errs.ErrArgs.WrapMsg("userID or ip cannot be empty")
		}
		limits[i] = &database.LimitUserLoginIP{
			UserID: limit.UserID,
			IP:     limit.Ip,
		}
	}

	// 3. ɾûIP¼
	if err := l.svcCtx.AdminDB.DelUserIPLimitLogin(l.ctx, limits); err != nil {
		l.Errorf("DelUserIPLimitLogin failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete user IP limit login")
	}

	l.Infof("Deleted user IP limit login: count=%d", len(req.Limits))
	return &admin.DelUserIPLimitLoginResp{}, nil
}
