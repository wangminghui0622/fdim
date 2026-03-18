package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/pkg/msgprocessor"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSeqMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSeqMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSeqMessageLogic {
	return &GetSeqMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSeqMessageLogic) GetSeqMessage(req *msg.GetSeqMessageReq) (*msg.GetSeqMessageResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.GetSeqMessageResp{
		Msgs:             make(map[string]*sdkws.PullMsgs),
		NotificationMsgs: make(map[string]*sdkws.PullMsgs),
	}

	for _, conv := range req.Conversations {
		if conv.ConversationID == "" || len(conv.Seqs) == 0 {
			continue
		}

		msgDocs, err := l.svcCtx.MsgCache.GetMessagesBySeq(l.ctx, conv.ConversationID, conv.Seqs)
		if err != nil {
			l.Errorw("GetMessagesBySeq failed", logx.Field("conversationID", conv.ConversationID), logx.Field("error", err))
			continue
		}
		if len(msgDocs) == 0 {
			continue
		}

		pullMsgs := &sdkws.PullMsgs{Msgs: make([]*sdkws.MsgData, 0, len(msgDocs))}
		for _, doc := range msgDocs {
			pullMsgs.Msgs = append(pullMsgs.Msgs, docToMsgData(doc))
		}

		if msgprocessor.IsNotification(conv.ConversationID) {
			resp.NotificationMsgs[conv.ConversationID] = pullMsgs
		} else {
			resp.Msgs[conv.ConversationID] = pullMsgs
		}
	}

	return resp, nil
}
