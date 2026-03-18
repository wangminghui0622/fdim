package logic

import (
	"context"

	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetSendMsgStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetSendMsgStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSendMsgStatusLogic {
	return &SetSendMsgStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetSendMsgStatusLogic) SetSendMsgStatus(req *msg.SetSendMsgStatusReq) (*msg.SetSendMsgStatusResp, error) {
	return &msg.SetSendMsgStatusResp{}, nil
}
