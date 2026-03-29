package logic

import (
	"context"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/bot"
	"fdim/rpc/bot/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentLogic {
	return &CreateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentLogic) CreateAgent(req *bot.CreateAgentReq) (*bot.CreateAgentResp, error) {
	// 1. ֤
	if req.Agent == nil {
		return nil, errs.ErrArgs.WrapMsg("agent cannot be nil")
	}
	if req.Agent.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user ID cannot be empty")
	}

	// 2. ôʱ
	createTime := req.Agent.CreateTime
	if createTime == 0 {
		createTime = time.Now().Unix()
	}

	// 3. תΪݿģ
	agent := &database.Agent{
		UserID:     req.Agent.UserID,
		Nickname:   req.Agent.Nickname,
		FaceURL:    req.Agent.FaceURL,
		URL:        req.Agent.Url,
		Key:        req.Agent.Key,
		Identity:   req.Agent.Identity,
		Model:      req.Agent.Model,
		Prompts:    req.Agent.Prompts,
		CreateTime: createTime,
	}

	// 4. Agent
	if err := l.svcCtx.BotDB.CreateAgent(l.ctx, agent); err != nil {
		l.Errorf("CreateAgent failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to create agent")
	}

	l.Infof("Agent created successfully: userID=%s", req.Agent.UserID)
	return &bot.CreateAgentResp{}, nil
}
