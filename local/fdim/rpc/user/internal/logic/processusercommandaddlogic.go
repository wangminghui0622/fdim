package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessUserCommandAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessUserCommandAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandAddLogic {
	return &ProcessUserCommandAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessUserCommandAddLogic) ProcessUserCommandAdd(req *user.ProcessUserCommandAddReq) (*user.ProcessUserCommandAddResp, error) {
	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// TODO: 实现 AddUserCommand 数据库方法
	// 当前返回未实现错误
	return nil, fmt.Errorf("ProcessUserCommandAdd not fully implemented yet: AddUserCommand database method needed")
}
