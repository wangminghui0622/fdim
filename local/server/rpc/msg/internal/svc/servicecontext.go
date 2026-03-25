package svc

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/pkg/mq"
	"fdim/pkg/notification"
	storagecache "fdim/pkg/storage/cache"
	"fdim/pkg/storage/cache/redis"
	"fdim/pkg/storage/controller"
	"fdim/pkg/storage/database/mgo"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/config"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	MongoDB *database.MongoDB
	Redis   *redis2.Client
	// 消息存储与缓存（与官方一致，使用 controller 层）
	MsgDatabase        controller.CommonMsgDatabase
	MsgCache           cache.MsgCache                // 使用 pkg/cache 接口（兼容旧代码）
	MsgDB              database.MsgDatabase          // 简化的数据库接口（兼容旧代码）
	SeqConversation    storagecache.SeqConversationCache // 官方 seq 分配器（SendMsg 必须使用此接口）
	// NATS producers
	ToRedisProducer mq.Producer // 发送到 toRedis topic，触发 msgtransfer 链路
	ToPushProducer  mq.Producer // 直接发送到 toPush topic，用于回退时推送
	// 内部 SendMsg 函数引用（用于发送通知，避免循环依赖）
	SendMsgFunc func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error)
	// 通知发送器（与官方一致）
	NotificationSender *notification.NotificationSender
	// RPC 客户端
	ConversationClient conversation.ConversationClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Host, c.Mongo.Database, c.Mongo.Username, c.Mongo.Password)

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化消息数据库（使用官方 storage 层接口）
	msgDocModel, err := mgo.NewMsgMongo(mongoDB.GetDatabase())
	if err != nil {
		logx.Errorf("Failed to create MsgMongo: %v", err)
	}
	
	// 初始化消息缓存（使用官方 storage 层接口）
	msgCacheModel := redis.NewMsgCache(redisClient, msgDocModel)
	
	// 初始化 SeqConversation（与官方一致）
	seqConversation, err := mgo.NewSeqConversationMongo(mongoDB.GetDatabase())
	if err != nil {
		logx.Errorf("Failed to create SeqConversation: %v", err)
	}
	seqConversationCache := redis.NewSeqConversationCacheRedis(redisClient, seqConversation)
	
	// 初始化 SeqUser（与官方一致）
	seqUser, err := mgo.NewSeqUserMongo(mongoDB.GetDatabase())
	if err != nil {
		logx.Errorf("Failed to create SeqUser: %v", err)
	}
	seqUserCache := redis.NewSeqUserCacheRedis(redisClient, seqUser)

	// 初始化 NATS Producers
	var toRedisProducer, toPushProducer mq.Producer
	
	// toRedis producer
	if len(c.Nats.Brokers) > 0 && c.Nats.Topics.ToRedis != "" {
		p, err := mq.NewNatsProducer(c.Nats.Brokers, c.Nats.Topics.ToRedis)
		if err != nil {
			logx.Errorf("Failed to create NATS toRedis producer: %v", err)
		} else {
			toRedisProducer = p
			logx.Infof("NATS toRedis producer initialized: brokers=%v, topic=%q", c.Nats.Brokers, c.Nats.Topics.ToRedis)
		}
	} else {
		logx.Errorf("NATS toRedis producer not initialized: Brokers=%v, ToRedis=%q", c.Nats.Brokers, c.Nats.Topics.ToRedis)
	}
	
	// toPush producer（用于回退时直接推送）
	if len(c.Nats.Brokers) > 0 && c.Nats.Topics.ToPush != "" {
		p, err := mq.NewNatsProducer(c.Nats.Brokers, c.Nats.Topics.ToPush)
		if err != nil {
			logx.Errorf("Failed to create NATS toPush producer: %v", err)
		} else {
			toPushProducer = p
			logx.Infof("NATS toPush producer initialized: brokers=%v, topic=%q", c.Nats.Brokers, c.Nats.Topics.ToPush)
		}
	}

	// 初始化 controller 层的 MsgDatabase（与官方一致）
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, msgCacheModel, seqUserCache, seqConversationCache, toRedisProducer)
	
	// 为了兼容旧代码，同时保留简化的 MsgCache 接口
	// 通过 SetSeqProvider 将 seq 操作委托给官方 SeqConversationCache，确保 seq 读写一致
	var simpleMsgCache cache.MsgCache
	if redisClient != nil {
		simpleMsgDocDB := database.NewMsgDocDatabase(mongoDB)
		rmc := cache.NewRedisMsgCache(redisClient, simpleMsgDocDB)
		if seqConversationCache != nil {
			rmc.SetSeqProvider(seqConversationCache)
			logx.Info("MsgCache SeqProvider set to official SeqConversationCache")
		}
		simpleMsgCache = rmc
	}

	// 初始化 Conversation RPC 客户端（增加超时和重试）
	var convClient conversation.ConversationClient
	// 支持 Etcd 服务发现或直连方式
	if len(c.ConversationRpc.Etcd.Hosts) > 0 || c.ConversationRpc.Target != "" || len(c.ConversationRpc.Endpoints) > 0 {
		// 设置更长的超时时间（10秒）
		if c.ConversationRpc.Timeout == 0 {
			c.ConversationRpc.Timeout = 10000 // 10秒
		}
		
		// 尝试连接，最多重试3次
		var conn zrpc.Client
		var err error
		for i := 0; i < 3; i++ {
			conn, err = zrpc.NewClient(c.ConversationRpc)
			if err == nil {
				convClient = conversation.NewConversationClient(conn.Conn())
				logx.Info("Conversation RPC client initialized successfully")
				break
			}
			logx.Errorf("Failed to create Conversation RPC client (attempt %d/3): %v", i+1, err)
			if i < 2 {
				time.Sleep(2 * time.Second) // 等待2秒后重试
			}
		}
		
		if convClient == nil {
			logx.Errorf("Failed to create Conversation RPC client after 3 attempts, will continue without it")
		}
	} else {
		logx.Error("Conversation RPC client not initialized: no Etcd, Target, or Endpoints configured")
	}

	svcCtx := &ServiceContext{
		Config:             c,
		MongoDB:            mongoDB,
		Redis:              redisClient,
		MsgDatabase:        msgDatabase,
		MsgCache:           simpleMsgCache,
		MsgDB:              database.NewMsgDocDatabase(mongoDB),
		SeqConversation:    seqConversationCache,
		ToRedisProducer:    toRedisProducer,
		ToPushProducer:     toPushProducer,
		ConversationClient: convClient,
		// NotificationSender will be initialized in msg.go after SendMsg is created
		NotificationSender: nil,
	}

	return svcCtx
}

// SetNotificationSender 设置通知发送器（在 msg.go 中调用）
func (s *ServiceContext) SetNotificationSender(sender *notification.NotificationSender) {
	s.NotificationSender = sender
}

// MsgCacheKey 辅助方法：生成消息在 Redis 中的缓存 key
func (s *ServiceContext) MsgCacheKey(conversationID string, seq int64) string {
	if s.Redis == nil {
		return ""
	}
	// 与 RedisMsgCache.getMessageKey 保持一致
	return fmt.Sprintf("msg:%s:%d", conversationID, seq)
}
