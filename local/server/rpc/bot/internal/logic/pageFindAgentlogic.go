package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/bot"
	"fdim/rpc/bot/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PageFindAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageFindAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageFindAgentLogic {
	return &PageFindAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageFindAgentLogic) PageFindAgent(req *bot.PageFindAgentReq) (*bot.PageFindAgentResp, error) {
	// 1. ҳAgent
	total, agents, err := l.svcCtx.BotDB.PageFindAgent(l.ctx, req.UserIDs, req.Pagination)
	if err != nil {
		l.Errorf("PageFindAgent failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find agents")
	}

	// 2. תΪprotobufʽ
	results := make([]*bot.Agent, 0, len(agents))
	for _, agent := range agents {
		results = append(results, &bot.Agent{
			UserID:     agent.UserID,
			Nickname:   agent.Nickname,
			FaceURL:    agent.FaceURL,
			Url:        agent.URL,
			Key:        agent.Key,
			Identity:   agent.Identity,
			Model:      agent.Model,
			Prompts:    agent.Prompts,
			CreateTime: agent.CreateTime,
		})
	}

	return &bot.PageFindAgentResp{
		Total:  total,
		Agents: results,
	}, nil
}
