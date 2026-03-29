package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessUserCommandDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessUserCommandDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandDeleteLogic {
	return &ProcessUserCommandDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessUserCommandDeleteLogic) ProcessUserCommandDelete(req *user.ProcessUserCommandDeleteReq) (*user.ProcessUserCommandDeleteResp, error) {
	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// TODO: 实现 DeleteUserCommand 数据库方?
	return nil, fmt.Errorf("ProcessUserCommandDelete not fully implemented yet: DeleteUserCommand database method needed")
}
