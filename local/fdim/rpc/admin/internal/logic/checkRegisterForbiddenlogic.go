package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckRegisterForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckRegisterForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckRegisterForbiddenLogic {
	return &CheckRegisterForbiddenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckRegisterForbiddenLogic) CheckRegisterForbidden(req *admin.CheckRegisterForbiddenReq) (*admin.CheckRegisterForbiddenResp, error) {
	// 1. 验证参数
	if req.Ip == "" {
		return nil, errs.ErrArgs.WrapMsg("ip cannot be empty")
	}

	// 2. 查找IP禁止记录
	forbiddens, err := l.svcCtx.AdminDB.FindIPForbidden(l.ctx, []string{req.Ip})
	if err != nil {
		l.Errorf("FindIPForbidden failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check IP forbidden")
	}

	// 3. 检查是否禁止注册
	for _, forbidden := range forbiddens {
		if forbidden.LimitRegister {
			return nil, fmt.Errorf("ip %s is forbidden for registration", req.Ip)
		}
	}

	return &admin.CheckRegisterForbiddenResp{}, nil
}
