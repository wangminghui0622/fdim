package logic

import (
	"context"

	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllowRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllowRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllowRegisterLogic {
	return &GetAllowRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllowRegisterLogic) GetAllowRegister(req *chat.GetAllowRegisterReq) (*chat.GetAllowRegisterResp, error) {
	return &chat.GetAllowRegisterResp{
		AllowRegister: l.svcCtx.Config.AllowRegister,
	}, nil
}
