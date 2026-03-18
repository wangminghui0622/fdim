package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessUserCommandGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessUserCommandGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandGetLogic {
	return &ProcessUserCommandGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessUserCommandGetLogic) ProcessUserCommandGet(req *user.ProcessUserCommandGetReq) (*user.ProcessUserCommandGetResp, error) {
	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// TODO: 实现 GetUserCommands 数据库方法
	return nil, fmt.Errorf("ProcessUserCommandGet not fully implemented yet: GetUserCommands database method needed")
}
