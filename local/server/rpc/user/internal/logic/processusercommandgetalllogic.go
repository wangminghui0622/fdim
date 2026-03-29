package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessUserCommandGetAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessUserCommandGetAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessUserCommandGetAllLogic {
	return &ProcessUserCommandGetAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessUserCommandGetAllLogic) ProcessUserCommandGetAll(req *user.ProcessUserCommandGetAllReq) (*user.ProcessUserCommandGetAllResp, error) {
	// 权限验证
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// TODO: 实现 GetAllUserCommands 数据库方?
	return nil, fmt.Errorf("ProcessUserCommandGetAll not fully implemented yet: GetAllUserCommands database method needed")
}
