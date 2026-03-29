package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetConversationMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetConversationMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationMaxSeqLogic {
	return &SetConversationMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetConversationMaxSeqLogic) SetConversationMaxSeq(req *conversation.SetConversationMaxSeqReq) (*conversation.SetConversationMaxSeqResp, error) {
	resp := &conversation.SetConversationMaxSeqResp{}

	if len(req.OwnerUserID) == 0 {
		return nil, errs.ErrArgs.WrapMsg("ownerUserID is required")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}

	// 1. 更新会话表中?max_seq
	args := map[string]interface{}{
		"max_seq": req.MaxSeq,
	}

	_, err := l.svcCtx.ConversationDB.UpdateByMap(l.ctx, req.OwnerUserID, req.ConversationID, args)
	if err != nil {
		return nil, fmt.Errorf("failed to set conversation max seq: %w", err)
	}

	// 2. 通知 Msg 模块更新缓存中的会话 MaxSeq
	if l.svcCtx.MsgClient != nil {
		_, err := l.svcCtx.MsgClient.SetUserConversationMaxSeq(l.ctx, &msg.SetUserConversationMaxSeqReq{
			ConversationID: req.ConversationID,
			OwnerUserID:    req.OwnerUserID,
			MaxSeq:         req.MaxSeq,
		})
		if err != nil {
			l.Errorf("failed to call SetUserConversationMaxSeq: %v", err)
		}
	}

	return resp, nil
}
