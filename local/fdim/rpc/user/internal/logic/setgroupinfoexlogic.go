package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupInfoExLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupInfoExLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupInfoExLogic {
	return &SetGroupInfoExLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetGroupInfoExLogic) SetGroupInfoEx(req *user.SetGroupInfoExReq) (*user.SetGroupInfoExResp, error) {
	// TODO: 实现设置群组信息（扩展）的逻辑
	return nil, fmt.Errorf("SetGroupInfoEx not fully implemented yet")
}
