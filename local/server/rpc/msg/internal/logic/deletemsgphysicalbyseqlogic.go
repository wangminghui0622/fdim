package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMsgPhysicalBySeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMsgPhysicalBySeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgPhysicalBySeqLogic {
	return &DeleteMsgPhysicalBySeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteMsgPhysicalBySeq 物理删除指定会话中指?seq 的消息?
func (l *DeleteMsgPhysicalBySeqLogic) DeleteMsgPhysicalBySeq(req *msg.DeleteMsgPhysicalBySeqReq) (*msg.DeleteMsgPhysicalBySeqResp, error) {
	if req.ConversationID == "" || len(req.Seqs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("conversationID and seqs are required")
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	// ?Mongo 中删?
	if err := l.svcCtx.MsgDB.DeleteMessagesBySeq(l.ctx, req.ConversationID, req.Seqs); err != nil {
		l.Errorw("DeleteMessagesBySeq failed", logx.Field("conversationID", req.ConversationID), logx.Field("seqs", req.Seqs), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to physically delete messages")
	}

	// ?Redis 清理缓存（忽略错误）
	if l.svcCtx.Redis != nil {
		for _, seq := range req.Seqs {
			key := l.svcCtx.MsgCacheKey(req.ConversationID, seq)
			if key == "" {
				continue
			}
			if err := l.svcCtx.Redis.Del(l.ctx, key).Err(); err != nil {
				l.Errorw("failed to delete message cache", logx.Field("key", key), logx.Field("error", err))
			}
		}
	}

	return &msg.DeleteMsgPhysicalBySeqResp{}, nil
}
