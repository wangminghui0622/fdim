package logic

import (
	"context"

	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendSimpleMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSimpleMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSimpleMsgLogic {
	return &SendSimpleMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSimpleMsgLogic) SendSimpleMsg(req *msg.SendSimpleMsgReq) (*msg.SendSimpleMsgResp, error) {
	// 复用 SendMsg 的实现逻辑
	sendLogic := NewSendMsgLogic(l.ctx, l.svcCtx)
	resp, err := sendLogic.SendMsg(&msg.SendMsgReq{
		MsgData: req.MsgData,
	})
	if err != nil {
		return nil, err
	}

	return &msg.SendSimpleMsgResp{
		ServerMsgID: resp.ServerMsgID,
		ClientMsgID: resp.ClientMsgID,
		SendTime:    resp.SendTime,
		Modify:      resp.Modify,
	}, nil
}
