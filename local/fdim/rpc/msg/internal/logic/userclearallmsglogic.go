package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserClearAllMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserClearAllMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClearAllMsgLogic {
	return &UserClearAllMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserClearAllMsgLogic) UserClearAllMsg(req *msg.UserClearAllMsgReq) (*msg.UserClearAllMsgResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	// 获取用户所有会话
	var conversationIDs []string
	if l.svcCtx.ConversationClient != nil {
		convResp, err := l.svcCtx.ConversationClient.GetConversationIDs(l.ctx, &conversation.GetConversationIDsReq{
			UserID: req.UserID,
		})
		if err != nil {
			l.Errorw("GetConversationIDs failed", logx.Field("error", err))
		} else {
			conversationIDs = convResp.ConversationIDs
		}
	}

	// 对每个会话设置 minSeq = maxSeq + 1（等同于清除所有消息）
	for _, convID := range conversationIDs {
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		if err := l.svcCtx.MsgCache.SetMinSeq(l.ctx, convID, maxSeq+1); err != nil {
			l.Errorw("SetMinSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
		}
	}

	return &msg.UserClearAllMsgResp{}, nil
}
