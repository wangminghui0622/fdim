package svc

import (
	"fdim/Infrastructure_service/cron/internal/config"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/third"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	MsgClient          msg.MsgClient
	ConversationClient conversation.ConversationClient
	ThirdClient        third.ThirdClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	var msgClient msg.MsgClient
	var conversationClient conversation.ConversationClient
	var thirdClient third.ThirdClient

	if c.MsgRpc.Etcd.Key != "" || c.MsgRpc.Target != "" {
		msgConn := zrpc.MustNewClient(c.MsgRpc).Conn()
		msgClient = msg.NewMsgClient(msgConn)
	}
	if c.ConversationRpc.Etcd.Key != "" || c.ConversationRpc.Target != "" {
		conversationConn := zrpc.MustNewClient(c.ConversationRpc).Conn()
		conversationClient = conversation.NewConversationClient(conversationConn)
	}
	if c.ThirdRpc.Etcd.Key != "" || c.ThirdRpc.Target != "" {
		thirdConn := zrpc.MustNewClient(c.ThirdRpc).Conn()
		thirdClient = third.NewThirdClient(thirdConn)
	}

	return &ServiceContext{
		Config:             c,
		MsgClient:          msgClient,
		ConversationClient: conversationClient,
		ThirdClient:        thirdClient,
	}
}
