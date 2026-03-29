package svc

import (
	"fdim/Infrastructure_service/push/internal/config"
	"fdim/Infrastructure_service/push/internal/offlinepush"
	"fdim/pkg/cache"
	"fdim/pkg/mq"
	"fdim/protocol/user"
	"fdim/protocol/msggateway"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config               config.Config
	PushConsumer         mq.Consumer
	OfflinePushConsumer  mq.Consumer
	MessageGatewayClient msggateway.MsgGatewayClient
	GroupClient          user.GroupClient
	OfflinePusher        offlinepush.OfflinePusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	// ʼ MessageGateway ͻ
	var messageGatewayClient msggateway.MsgGatewayClient
	if c.MessageGatewayRpc.Etcd.Key != "" || c.MessageGatewayRpc.Target != "" {
		messageGatewayConn := zrpc.MustNewClient(c.MessageGatewayRpc).Conn()
		messageGatewayClient = msggateway.NewMsgGatewayClient(messageGatewayConn)
	}

	// ʼ Redis ͻˣ token 洢
	redisClient := cache.NewRedisClient(c.Cache)

	// ʼ Group ͻ
	var groupClient user.GroupClient
	if c.GroupRpc.Etcd.Key != "" || c.GroupRpc.Target != "" {
		groupConn := zrpc.MustNewClient(c.GroupRpc).Conn()
		groupClient = user.NewGroupClient(groupConn)
	}

	// ʼ NATS Consumer
	var pushConsumer mq.Consumer
	if c.Nats.ToPushTopic == "" || c.Nats.ToPushGroupID == "" {
		logx.Errorf("push consumer not initialized: empty NATS config (ToPushTopic=%q, ToPushGroupID=%q)", c.Nats.ToPushTopic, c.Nats.ToPushGroupID)
	} else {
		consumer, err := mq.NewNatsConsumer(c.Nats.Brokers, c.Nats.ToPushGroupID, c.Nats.ToPushTopic)
		if err != nil {
			logx.Errorf("push consumer not initialized: failed to create NATS consumer (topic=%q, groupID=%q): %v", c.Nats.ToPushTopic, c.Nats.ToPushGroupID, err)
		} else {
			pushConsumer = consumer
		}
	}

	var offlinePushConsumer mq.Consumer
	if c.Nats.ToOfflinePushTopic == "" || c.Nats.ToOfflinePushGroupID == "" {
		logx.Errorf("offline push consumer not initialized: empty NATS config (ToOfflinePushTopic=%q, ToOfflinePushGroupID=%q)", c.Nats.ToOfflinePushTopic, c.Nats.ToOfflinePushGroupID)
	} else {
		consumer, err := mq.NewNatsConsumer(c.Nats.Brokers, c.Nats.ToOfflinePushGroupID, c.Nats.ToOfflinePushTopic)
		if err != nil {
			logx.Errorf("offline push consumer not initialized: failed to create NATS consumer (topic=%q, groupID=%q): %v", c.Nats.ToOfflinePushTopic, c.Nats.ToOfflinePushGroupID, err)
		} else {
			offlinePushConsumer = consumer
		}
	}

	// ʼ
	pushConfig := offlinepush.PushConfig{
		Enable:            c.Push.Enable,
		FcmServerKey:      c.Push.FcmServerKey,
		JPushAppKey:       c.Push.JPushAppKey,
		JPushMasterSecret: c.Push.JPushMasterSecret,
	}
	offlinePusher, _ := offlinepush.NewOfflinePusherWithConfig(pushConfig, redisClient)

	return &ServiceContext{
		Config:               c,
		PushConsumer:         pushConsumer,
		OfflinePushConsumer:  offlinePushConsumer,
		MessageGatewayClient: messageGatewayClient,
		GroupClient:          groupClient,
		OfflinePusher:        offlinePusher,
	}
}
