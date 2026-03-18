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

type SetConversationMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetConversationMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConversationMinSeqLogic {
	return &SetConversationMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetConversationMinSeqLogic) SetConversationMinSeq(req *conversation.SetConversationMinSeqReq) (*conversation.SetConversationMinSeqResp, error) {
	resp := &conversation.SetConversationMinSeqResp{}

	if len(req.OwnerUserID) == 0 {
		return nil, errs.ErrArgs.WrapMsg("ownerUserID is required")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}

	// 1. 更新会话表中的 min_seq
	args := map[string]interface{}{
		"min_seq": req.MinSeq,
	}

	_, err := l.svcCtx.ConversationDB.UpdateByMap(l.ctx, req.OwnerUserID, req.ConversationID, args)
	if err != nil {
		return nil, fmt.Errorf("failed to set conversation min seq: %w", err)
	}

	// 2. 通知 Msg 模块更新缓存中的会话 MinSeq
	if l.svcCtx.MsgClient != nil {
		_, err := l.svcCtx.MsgClient.SetUserConversationMinSeq(l.ctx, &msg.SetUserConversationMinSeqReq{
			ConversationID: req.ConversationID,
			OwnerUserID:    req.OwnerUserID,
			MinSeq:         req.MinSeq,
		})
		if err != nil {
			l.Errorf("failed to call SetUserConversationMinSeq: %v", err)
		}
	}

	return resp, nil
}
