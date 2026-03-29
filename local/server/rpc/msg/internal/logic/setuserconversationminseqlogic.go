package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserConversationMinSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserConversationMinSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserConversationMinSeqLogic {
	return &SetUserConversationMinSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserConversationMinSeqLogic) SetUserConversationMinSeq(req *msg.SetUserConversationMinSeqReq) (*msg.SetUserConversationMinSeqResp, error) {
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}
	if len(req.OwnerUserID) == 0 {
		return nil, errs.ErrArgs.WrapMsg("ownerUserID is required")
	}
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	//为每个 ownerUserID 设置该会话的最小已保留 seq（一般对应清理历史消息的位置?
	for _, owner := range req.OwnerUserID {
		if owner == "" {
			continue
		}
		// 这里只在缓存中记录最?seq，具体消息物理删除由其他任务负责
		if err := l.svcCtx.MsgCache.SetMinSeq(l.ctx, req.ConversationID, req.MinSeq); err != nil {
			l.Errorw("SetMinSeq failed", logx.Field("conversationID", req.ConversationID), logx.Field("ownerUserID", owner), logx.Field("minSeq", req.MinSeq), logx.Field("error", err))
			return nil, errs.ErrInternalServer.WrapMsg("failed to set conversation min seq")
		}
	}

	return &msg.SetUserConversationMinSeqResp{}, nil
}
