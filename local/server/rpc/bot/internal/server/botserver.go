package server

import (
	"context"

	"fdim/protocol/bot"
	"fdim/rpc/bot/internal/logic"
	"fdim/rpc/bot/internal/svc"
)

type BotServer struct {
	bot.UnimplementedBotServer
	svcCtx *svc.ServiceContext
}

func NewBotServer(svcCtx *svc.ServiceContext) *BotServer {
	return &BotServer{
		svcCtx: svcCtx,
	}
}

func (s *BotServer) CreateAgent(ctx context.Context, req *bot.CreateAgentReq) (*bot.CreateAgentResp, error) {
	l := logic.NewCreateAgentLogic(ctx, s.svcCtx)
	return l.CreateAgent(req)
}

func (s *BotServer) UpdateAgent(ctx context.Context, req *bot.UpdateAgentReq) (*bot.UpdateAgentResp, error) {
	l := logic.NewUpdateAgentLogic(ctx, s.svcCtx)
	return l.UpdateAgent(req)
}

func (s *BotServer) PageFindAgent(ctx context.Context, req *bot.PageFindAgentReq) (*bot.PageFindAgentResp, error) {
	l := logic.NewPageFindAgentLogic(ctx, s.svcCtx)
	return l.PageFindAgent(req)
}

func (s *BotServer) DeleteAgent(ctx context.Context, req *bot.DeleteAgentReq) (*bot.DeleteAgentResp, error) {
	l := logic.NewDeleteAgentLogic(ctx, s.svcCtx)
	return l.DeleteAgent(req)
}

func (s *BotServer) SendBotMessage(ctx context.Context, req *bot.SendBotMessageReq) (*bot.SendBotMessageResp, error) {
	l := logic.NewSendBotMessageLogic(ctx, s.svcCtx)
	return l.SendBotMessage(req)
}
