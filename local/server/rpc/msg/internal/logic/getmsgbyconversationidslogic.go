package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMsgByConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMsgByConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMsgByConversationIDsLogic {
	return &GetMsgByConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMsgByConversationIDsLogic) GetMsgByConversationIDs(req *msg.GetMsgByConversationIDsReq) (*msg.GetMsgByConversationIDsResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &msg.GetMsgByConversationIDsResp{
		MsgDatas: make(map[string]*sdkws.MsgData),
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
		resp.MsgDatas[convID] = docToMsgData(msgDocs[0])
	}

	return resp, nil
}
