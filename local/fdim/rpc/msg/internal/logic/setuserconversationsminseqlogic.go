package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserConversationsMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserConversationsMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserConversationsMinSeqLogic {
	return &SetUserConversationsMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserConversationsMinSeqLogic) SetUserConversationsMinSeq(req *msg.SetUserConversationsMinSeqReq) (*msg.SetUserConversationsMinSeqResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}

	if err := l.svcCtx.MsgCache.SetMinSeq(l.ctx, req.ConversationID, req.Seq); err != nil {
		return nil, errs.WrapMsg(err, "SetMinSeq failed")
	}

	return &msg.SetUserConversationsMinSeqResp{}, nil
}
