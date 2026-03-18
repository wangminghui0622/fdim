package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessUserCommandUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessUserCommandUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandUpdateLogic {
	return &ProcessUserCommandUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessUserCommandUpdateLogic) ProcessUserCommandUpdate(req *user.ProcessUserCommandUpdateReq) (*user.ProcessUserCommandUpdateResp, error) {
	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// TODO: 实现 UpdateUserCommand 数据库方法
	return nil, fmt.Errorf("ProcessUserCommandUpdate not fully implemented yet: UpdateUserCommand database method needed")
}
