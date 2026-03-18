package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMsgPhysicalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMsgPhysicalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgPhysicalLogic {
	return &DeleteMsgPhysicalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteMsgPhysical 按会话列表 + 时间戳，物理删除历史消息。
// 这里直接调用 MsgDB.DeleteMessagesByTimeBefore 做一次性清理，未对 Redis 做细粒度同步。
func (l *DeleteMsgPhysicalLogic) DeleteMsgPhysical(req *msg.DeleteMsgPhysicalReq) (*msg.DeleteMsgPhysicalResp, error) {
	if len(req.ConversationIDs) == 0 || req.Timestamp == 0 {
		return nil, errs.ErrArgs.WrapMsg("conversationIDs and timestamp are required")
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	if err := l.svcCtx.MsgDB.DeleteMessagesByTimeBefore(l.ctx, req.ConversationIDs, req.Timestamp); err != nil {
		l.Errorw("DeleteMessagesByTimeBefore failed", logx.Field("conversationIDs", req.ConversationIDs), logx.Field("timestamp", req.Timestamp), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to physically delete messages by time")
	}

	// Redis 中的历史缓存不做精确删除，交给过期时间与后续重建逻辑处理。
	return &msg.DeleteMsgPhysicalResp{}, nil
}
