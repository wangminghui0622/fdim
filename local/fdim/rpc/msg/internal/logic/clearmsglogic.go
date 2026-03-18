package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClearMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearMsgLogic {
	return &ClearMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearMsgLogic) ClearMsg(req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	if len(req.Conversations) == 0 {
		return &msg.ClearMsgResp{}, nil
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	for _, conv := range req.Conversations {
		if conv == nil || conv.ConversationID == "" {
			continue
		}

		// 获取当前最大 seq
		maxSeq, err := l.svcCtx.MsgDB.GetMaxSeq(l.ctx, conv.ConversationID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", conv.ConversationID), logx.Field("error", err))
			continue
		}
		if maxSeq == 0 {
			continue
		}

		// 删除该会话的所有消息
		seqs := make([]int64, 0, maxSeq)
		for s := int64(1); s <= maxSeq; s++ {
			seqs = append(seqs, s)
		}
		if err := l.svcCtx.MsgDB.DeleteMessagesBySeq(l.ctx, conv.ConversationID, seqs); err != nil {
			l.Errorw("DeleteMessagesBySeq failed", logx.Field("conversationID", conv.ConversationID), logx.Field("error", err))
		}
	}

	return &msg.ClearMsgResp{}, nil
}
