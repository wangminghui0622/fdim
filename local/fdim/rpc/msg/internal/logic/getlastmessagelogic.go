package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMessageLogic {
	return &GetLastMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLastMessageLogic) GetLastMessage(req *msg.GetLastMessageReq) (*msg.GetLastMessageResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.GetLastMessageResp{
		Msgs: make(map[string]*sdkws.MsgData),
	}

	for _, convID := range req.ConversationIDs {
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil || maxSeq <= 0 {
			continue
		}
		msgDocs, err := l.svcCtx.MsgCache.GetMessagesBySeq(l.ctx, convID, []int64{maxSeq})
		if err != nil || len(msgDocs) == 0 {
			continue
		}
		resp.Msgs[convID] = docToMsgData(msgDocs[0])
	}

	return resp, nil
}
