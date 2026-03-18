package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastMessageSeqByTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastMessageSeqByTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMessageSeqByTimeLogic {
	return &GetLastMessageSeqByTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLastMessageSeqByTimeLogic) GetLastMessageSeqByTime(req *msg.GetLastMessageSeqByTimeReq) (*msg.GetLastMessageSeqByTimeResp, error) {
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}
	if req.Time <= 0 {
		return nil, errs.ErrArgs.WrapMsg("time must be positive")
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	seq, err := l.svcCtx.MsgDB.GetLastMessageSeqByTime(l.ctx, req.ConversationID, req.Time)
	if err != nil {
		l.Errorw("GetLastMessageSeqByTime failed", logx.Field("conversationID", req.ConversationID), logx.Field("time", req.Time), logx.Field("error", err))
		return nil, errs.ErrInternalServer.WrapMsg("failed to query last message seq by time")
	}

	return &msg.GetLastMessageSeqByTimeResp{
		Seq: seq,
	}, nil
}
