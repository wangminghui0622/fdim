package logic

import (
	"context"

	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetAllowRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetAllowRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAllowRegisterLogic {
	return &SetAllowRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetAllowRegisterLogic) SetAllowRegister(req *chat.SetAllowRegisterReq) (*chat.SetAllowRegisterResp, error) {
	// 简化实现：直接更新配置
	// 实际应该保存到数据库或配置文件
	l.svcCtx.Config.AllowRegister = req.AllowRegister

	l.Infof("Set allow register: %v", req.AllowRegister)
	return &chat.SetAllowRegisterResp{}, nil
}
