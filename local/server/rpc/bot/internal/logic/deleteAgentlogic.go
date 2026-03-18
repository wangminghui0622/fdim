package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/bot"
	"fdim/rpc/bot/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentLogic {
	return &DeleteAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAgentLogic) DeleteAgent(req *bot.DeleteAgentReq) (*bot.DeleteAgentResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user IDs cannot be empty")
	}

	// 2. 删除Agent
	if err := l.svcCtx.BotDB.DeleteAgent(l.ctx, req.UserIDs); err != nil {
		l.Errorf("DeleteAgent failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete agents")
	}

	l.Infof("Deleted agents: count=%d", len(req.UserIDs))
	return &bot.DeleteAgentResp{}, nil
}
