package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DestructMsgsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDestructMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DestructMsgsLogic {
	return &DestructMsgsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DestructMsgsLogic) DestructMsgs(req *msg.DestructMsgsReq) (*msg.DestructMsgsResp, error) {
	if req.Timestamp == 0 {
		return nil, errs.ErrArgs.WrapMsg("timestamp is required")
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	// 简化实现：全局按时间清理，而不是分会话 + limit 精确控制
	// 由于 MsgDatabase 当前没有列出所?conversationID 的方法，这里直接用空会话列表留给后续扩展?
	if err := l.svcCtx.MsgDB.DeleteMessagesByTimeBefore(l.ctx, []string{}, req.Timestamp); err != nil {
		l.Errorw("DeleteMessagesByTimeBefore failed", logx.Field("timestamp", req.Timestamp), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to destruct messages")
	}

	// count 暂时返回 0（未做精确统计）
	return &msg.DestructMsgsResp{
		Count: 0,
	}, nil
}
