package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/bot"
	"fdim/rpc/bot/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentLogic {
	return &UpdateAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAgentLogic) UpdateAgent(req *bot.UpdateAgentReq) (*bot.UpdateAgentResp, error) {
	// 1. ֤
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user ID cannot be empty")
	}

	// 2. ֶ
	update := make(map[string]interface{})
	if req.Nickname != nil && *req.Nickname != "" {
		update["nickname"] = *req.Nickname
	}
	if req.FaceURL != nil && *req.FaceURL != "" {
		update["face_url"] = *req.FaceURL
	}
	if req.Url != nil && *req.Url != "" {
		update["url"] = *req.Url
	}
	if req.Key != nil && *req.Key != "" {
		update["key"] = *req.Key
	}
	if req.Identity != nil && *req.Identity != "" {
		update["identity"] = *req.Identity
	}
	if req.Model != nil && *req.Model != "" {
		update["model"] = *req.Model
	}
	if req.Prompts != nil && *req.Prompts != "" {
		update["prompts"] = *req.Prompts
	}

	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update fields provided")
	}

	// 3. Agent
	if err := l.svcCtx.BotDB.UpdateAgent(l.ctx, req.UserID, update); err != nil {
		l.Errorf("UpdateAgent failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to update agent")
	}

	l.Infof("Agent updated successfully: userID=%s", req.UserID)
	return &bot.UpdateAgentResp{}, nil
}
