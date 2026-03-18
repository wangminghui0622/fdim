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

type AddUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserIPLimitLoginLogic {
	return &AddUserIPLimitLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserIPLimitLoginLogic) AddUserIPLimitLogin(req *admin.AddUserIPLimitLoginReq) (*admin.AddUserIPLimitLoginResp, error) {
	// 1. 验证参数
	if len(req.Limits) == 0 {
		return nil, errs.ErrArgs.WrapMsg("limits cannot be empty")
	}

	// 2. 构建限制列表
	limits := make([]*database.LimitUserLoginIP, len(req.Limits))
	now := time.Now()
	for i, limit := range req.Limits {
		if limit.UserID == "" || limit.Ip == "" {
			return nil, errs.ErrArgs.WrapMsg("userID or ip cannot be empty")
		}
		limits[i] = &database.LimitUserLoginIP{
			UserID:     limit.UserID,
			IP:         limit.Ip,
			CreateTime: now,
		}
	}

	// 3. 添加用户IP登录限制
	if err := l.svcCtx.AdminDB.AddUserIPLimitLogin(l.ctx, limits); err != nil {
		l.Errorf("AddUserIPLimitLogin failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add user IP limit login")
	}

	l.Infof("Added user IP limit login: count=%d", len(req.Limits))
	return &admin.AddUserIPLimitLoginResp{}, nil
}
