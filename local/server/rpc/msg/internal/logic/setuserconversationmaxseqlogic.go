package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserConversationMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserConversationMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserConversationMaxSeqLogic {
	return &SetUserConversationMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserConversationMaxSeqLogic) SetUserConversationMaxSeq(req *msg.SetUserConversationMaxSeqReq) (*msg.SetUserConversationMaxSeqResp, error) {
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}
	if len(req.OwnerUserID) == 0 {
		return nil, errs.ErrArgs.WrapMsg("ownerUserID is required")
	}
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	// 对于每个 ownerUserID，将会话的 MaxSeq 更新为给定值
	for _, owner := range req.OwnerUserID {
		if owner == "" {
			continue
		}
		if err := l.svcCtx.MsgCache.SetMaxSeq(l.ctx, req.ConversationID, req.MaxSeq); err != nil {
			l.Errorw("SetMaxSeq failed", logx.Field("conversationID", req.ConversationID), logx.Field("ownerUserID", owner), logx.Field("maxSeq", req.MaxSeq), logx.Field("error", err))
			return nil, errs.ErrInternalServer.WrapMsg("failed to set conversation max seq")
		}
	}

	return &msg.SetUserConversationMaxSeqResp{}, nil
}
