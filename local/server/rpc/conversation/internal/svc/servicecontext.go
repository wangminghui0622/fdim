package svc

import (
	"context"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	pkgconfig "fdim/pkg/config"
	"fdim/pkg/notification"
	"fdim/protocol/msg"
	"fdim/protocol/user"
	"fdim/rpc/conversation/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                       config.Config
	MongoDB                      *database.MongoDB
	Redis                        *redis.Client
	ConversationDB               *database.ConversationDatabase
	MsgClient                    msg.MsgClient
	UserClient                   user.UserClient
	GroupClient                  user.GroupClient
	ConversationNotification     *notification.ConversationNotificationSender
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始?MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Url, c.Mongo.Db, "", "")

	// 初始?Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始?ConversationDatabase
	conversationDB := database.NewConversationDatabase(mongoDB)

	// 初始?Msg RPC 客户端（使用 NewClient 而不?MustNewClient，避免循环依赖启动失败）
	var msgClient msg.MsgClient
	if c.MsgRpc.Etcd.Key != "" || c.MsgRpc.Target != "" {
		// 设置超时时间
		if c.MsgRpc.Timeout == 0 {
			c.MsgRpc.Timeout = 10000 // 10?
		}
		
		msgConn, err := zrpc.NewClient(c.MsgRpc)
		if err != nil {
			logx.Errorf("Failed to create Msg RPC client: %v, conversation notifications will be disabled", err)
		} else {
			msgClient = msg.NewMsgClient(msgConn.Conn())
			logx.Info("Msg RPC client initialized successfully")
		}
	}

	// 初始?User RPC 客户端（使用 NewClient 避免启动失败?
	var userClient user.UserClient
	var groupClient user.GroupClient
	if c.UserRpc.Etcd.Key != "" || c.UserRpc.Target != "" {
		// 设置超时时间
		if c.UserRpc.Timeout == 0 {
			c.UserRpc.Timeout = 10000 // 10?
		}
		
		userConn, err := zrpc.NewClient(c.UserRpc)
		if err != nil {
			logx.Errorf("Failed to create User RPC client: %v, user/group info will not be available", err)
		} else {
			userClient = user.NewUserClient(userConn.Conn())
			groupClient = user.NewGroupClient(userConn.Conn())
			logx.Info("User RPC client initialized successfully")
		}
	} else {
		logx.Infof("UserRpc not configured, user/group info will not be available")
	}

	// 初始化会话通知发送器
	notifConf := pkgconfig.DefaultNotification()
	var notifOpts []notification.NotificationSenderOptions
	if msgClient != nil {
		notifOpts = append(notifOpts, notification.WithRpcClient(func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
			return msgClient.SendMsg(ctx, req)
		}))
	} else {
		logx.Infof("MsgRpc not configured, conversation notifications will be disabled")
	}
	convNotification := notification.NewConversationNotificationSender(notifConf, notifOpts...)

	return &ServiceContext{
		Config:                   c,
		MongoDB:                  mongoDB,
		Redis:                    redisClient,
		ConversationDB:           conversationDB,
		MsgClient:                msgClient,
		UserClient:               userClient,
		GroupClient:              groupClient,
		ConversationNotification: convNotification,
	}
}
