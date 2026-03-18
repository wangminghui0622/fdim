package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationMaxSeqLogic {
	return &GetConversationMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationMaxSeqLogic) GetConversationMaxSeq(req *msg.GetConversationMaxSeqReq) (*msg.GetConversationMaxSeqResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, req.ConversationID)
	if err != nil {
		return nil, errs.WrapMsg(err, "GetMaxSeq failed")
	}

	return &msg.GetConversationMaxSeqResp{
		MaxSeq: maxSeq,
	}, nil
}
