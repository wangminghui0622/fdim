package logic

import (
	"context"

	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSendMsgStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSendMsgStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSendMsgStatusLogic {
	return &GetSendMsgStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSendMsgStatusLogic) GetSendMsgStatus(req *msg.GetSendMsgStatusReq) (*msg.GetSendMsgStatusResp, error) {
	return &msg.GetSendMsgStatusResp{
		Status: 0,
	}, nil
}
