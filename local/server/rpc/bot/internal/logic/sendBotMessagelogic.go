package logic

import (
	"context"
	"strings"

	"fdim/pkg/errs"
	"fdim/pkg/util/conversationutil"
	"fdim/protocol/bot"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/bot/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendBotMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendBotMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBotMessageLogic {
	return &SendBotMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendBotMessageLogic) SendBotMessage(req *bot.SendBotMessageReq) (*bot.SendBotMessageResp, error) {
	// 1. 校验参数
	if req.AgentID == "" {
		return nil, errs.ErrArgs.WrapMsg("agentID cannot be empty")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID cannot be empty")
	}

	// 2. 查询 Agent 信息
	_, agents, err := l.svcCtx.BotDB.PageFindAgent(l.ctx, []string{req.AgentID}, nil)
	if err != nil || len(agents) == 0 {
		return nil, errs.ErrArgs.WrapMsg("agent not found")
	}
	agent := agents[0]

	// 3. 从 conversationID 推导单聊 recvID / 群聊 groupID
	var recvID string
	var groupID string
	sessionType := int32(1)
	if conversationutil.IsGroupConversationID(req.ConversationID) {
		sessionType = 3
		groupID = strings.TrimPrefix(req.ConversationID, "sg_")
		if groupID == "" {
			return nil, errs.ErrArgs.WrapMsg("invalid group conversationID")
		}
	} else {
		// single conversation: si_userA_userB
		if !strings.HasPrefix(req.ConversationID, "si_") {
			return nil, errs.ErrArgs.WrapMsg("invalid conversationID")
		}
		parts := strings.Split(strings.TrimPrefix(req.ConversationID, "si_"), "_")
		if len(parts) != 2 {
			return nil, errs.ErrArgs.WrapMsg("invalid single conversationID")
		}
		// prefer the other side as recvID
		if parts[0] == agent.UserID {
			recvID = parts[1]
		} else if parts[1] == agent.UserID {
			recvID = parts[0]
		} else {
			recvID = parts[1]
		}
	}

	contentBytes := []byte(req.Content)

	msgData := &sdkws.MsgData{
		SendID:         agent.UserID,
		RecvID:         recvID,
		GroupID:        groupID,
		ClientMsgID:    req.Key,
		SenderNickname: agent.Nickname,
		SenderFaceURL:  agent.FaceURL,
		SessionType:    sessionType,
		ContentType:    req.ContentType,
		Content:        contentBytes,
		AttachedInfo:   req.Ex,
	}

	// 4. 调用 Msg RPC 发送消息
	if l.svcCtx.MsgRpc == nil {
		return nil, errs.ErrInternalServer.WrapMsg("Msg RPC client not initialized")
	}
	_, err = l.svcCtx.MsgRpc.SendMsg(l.ctx, &msg.SendMsgReq{
		MsgData: msgData,
	})
	if err != nil {
		l.Errorf("SendMsg via Msg RPC failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to send bot message")
	}

	l.Infof("Bot message sent: agentID=%s conversationID=%s", req.AgentID, req.ConversationID)
	return &bot.SendBotMessageResp{}, nil
}
