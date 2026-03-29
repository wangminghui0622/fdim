package svc

import (
	"context"

	"fdim/pkg/cache"
	pkgconfig "fdim/pkg/config"
	"fdim/pkg/database"
	"fdim/pkg/notification"
	"fdim/pkg/webhook"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/user/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	UserDB             *database.UserDatabase
	FriendDB           *database.FriendDatabase
	GroupDB            *database.GroupDatabase
	BlackDB            *database.BlackDatabase
	UserCache          *cache.UserCache
	FriendCache        *cache.FriendCache
	GroupCache         *cache.GroupCache
	FriendNotification *notification.FriendNotificationSender
	GroupNotification  *notification.GroupNotificationSender
	UserNotification   *notification.UserNotificationSender
	WebhookClient      *webhook.Client
	ConversationClient conversation.ConversationClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始?MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Host, c.Mongo.Database, c.Mongo.Username, c.Mongo.Password)

	// 初始?Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化数据库?
	userDB := database.NewUserDatabase(mongoDB)
	friendDB := database.NewFriendDatabase(mongoDB)
	groupDB := database.NewGroupDatabase(mongoDB)
	blackDB := database.NewBlackDatabase(mongoDB)

	// 初始化缓存层
	userCache := cache.NewUserCache(redisClient)
	friendCache := cache.NewFriendCache(redisClient)
	groupCache := cache.NewGroupCache(redisClient)

	// 初始化通知发送器
	var notificationOpts []notification.NotificationSenderOptions
	if c.MsgRpc.Etcd.Key != "" {
		msgRpcClient, err := zrpc.NewClient(c.MsgRpc)
		if err != nil {
			logx.Errorf("failed to connect msg.rpc: %v, notifications will be disabled", err)
		} else {
			msgClient := msg.NewMsgClient(msgRpcClient.Conn())
			notificationOpts = append(notificationOpts, notification.WithRpcClient(
				func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
					return msgClient.SendMsg(ctx, req)
				},
			))
			logx.Info("msg.rpc client connected, notifications enabled")
		}
	} else {
		logx.Infof("MsgRpc not configured, notifications will be disabled")
	}
	notifConf := pkgconfig.DefaultNotification()
	friendNotification := notification.NewFriendNotificationSender(notifConf, notificationOpts...)
	// 设置好友数据库，用于获取好友申请信息（包含reqMsg?
	friendNotification.SetFriendDB(friendDB)
	groupNotification := notification.NewGroupNotificationSender(notifConf, notificationOpts...)
	userNotification := notification.NewUserNotificationSender(notifConf, notificationOpts...)

	// 初始?Webhook 客户?
	webhookClient := webhook.NewWebhookClient(c.Webhook.URL)

	// 初始?Conversation RPC 客户?
	var convClient conversation.ConversationClient
	if c.ConversationRpc.Target != "" || len(c.ConversationRpc.Endpoints) > 0 {
		conn, err := zrpc.NewClient(c.ConversationRpc)
		if err != nil {
			logx.Errorf("Failed to create Conversation RPC client: %v", err)
		} else {
			convClient = conversation.NewConversationClient(conn.Conn())
			logx.Info("Conversation RPC client initialized in user service")
		}
	}

	return &ServiceContext{
		Config:             c,
		UserDB:             userDB,
		FriendDB:           friendDB,
		GroupDB:            groupDB,
		BlackDB:            blackDB,
		UserCache:          userCache,
		FriendCache:        friendCache,
		GroupCache:         groupCache,
		FriendNotification: friendNotification,
		GroupNotification:  groupNotification,
		UserNotification:   userNotification,
		WebhookClient:      webhookClient,
		ConversationClient: convClient,
	}
}
