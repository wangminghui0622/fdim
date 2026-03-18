package handler

import (
	"context"
	"fdim/Infrastructure_service/msgtransfer/internal/config"
	"fmt"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/pkg/mq"
	"fdim/protocol/conversation"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

// MsgTransferHandler 消息传输处理器
type MsgTransferHandler struct {
	config             *config.Config
	mongoDB            *database.MongoDB
	redisClient        *redis.Client
	msgCache           cache.MsgCache
	msgDB              database.MsgDatabase
	conversationClient conversation.ConversationClient
	toRedisConsumer    mq.Consumer
	toMongoConsumer    mq.Consumer
	toMongoProducer    mq.Producer
	toPushProducer     mq.Producer
}

// NewMsgTransferHandler 创建消息传输处理器
func NewMsgTransferHandler(ctx context.Context, cfg *config.Config) (*MsgTransferHandler, error) {
	// 初始化 MongoDB
	var mongoDB *database.MongoDB
	if cfg.MongoDB.URI != "" {
		mongoDB = database.NewMongoDB(
			cfg.MongoDB.URI,
			cfg.MongoDB.Database,
			cfg.MongoDB.Username,
			cfg.MongoDB.Password,
		)
	}

	// 初始化 Redis
	var redisClient *redis.Client
	if len(cfg.Cache) > 0 {
		logx.Infof("Initializing Redis client: Host=%s, Pass=%s", cfg.Cache[0].Host, cfg.Cache[0].Pass)
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Cache[0].Host,
			Password: cfg.Cache[0].Pass,
			DB:       0,
		})
		// 测试 Redis 连接
		if err := redisClient.Ping(ctx).Err(); err != nil {
			logx.Errorf("Redis connection test failed: %v", err)
		} else {
			logx.Info("Redis connection test successful")
		}
	}

	// 初始化消息数据库
	var msgDB database.MsgDatabase
	if mongoDB != nil {
		msgDB = database.NewMsgDocDatabase(mongoDB)
	}

	// 初始化消息缓存（需要 msgDB 作为参数）
	var msgCache *cache.RedisMsgCache
	if redisClient != nil && msgDB != nil {
		msgCache = cache.NewRedisMsgCache(redisClient, msgDB)
	}

	// 初始化 NATS Consumer 和 Producer
	toRedisConsumer, err := mq.NewNatsConsumer(cfg.Nats.Brokers, cfg.Nats.ToRedisGroupID, cfg.Nats.ToRedisTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to create toRedis consumer: %w", err)
	}

	toMongoConsumer, err := mq.NewNatsConsumer(cfg.Nats.Brokers, cfg.Nats.ToMongoGroupID, cfg.Nats.ToMongoTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to create toMongo consumer: %w", err)
	}

	toMongoProducer, err := mq.NewNatsProducer(cfg.Nats.Brokers, cfg.Nats.ToMongoTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to create toMongo producer: %w", err)
	}

	toPushProducer, err := mq.NewNatsProducer(cfg.Nats.Brokers, cfg.Nats.ToPushTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to create toPush producer: %w", err)
	}

	// 初始化 Conversation RPC 客户端（用于创建会话）
	var conversationClient conversation.ConversationClient
	if len(cfg.ConversationRpc.Etcd.Hosts) > 0 || cfg.ConversationRpc.Target != "" {
		// 设置超时时间
		if cfg.ConversationRpc.Timeout == 0 {
			cfg.ConversationRpc.Timeout = 10000 // 10秒
		}
		
		convConn, err := zrpc.NewClient(cfg.ConversationRpc)
		if err != nil {
			logx.Errorf("Failed to create Conversation RPC client: %v, conversation creation will be disabled", err)
		} else {
			conversationClient = conversation.NewConversationClient(convConn.Conn())
			logx.Info("Conversation RPC client initialized successfully")
		}
	} else {
		logx.Error("Conversation RPC not configured, conversation creation will be disabled")
	}

	return &MsgTransferHandler{
		config:             cfg,
		mongoDB:            mongoDB,
		redisClient:        redisClient,
		msgCache:           msgCache,
		msgDB:              msgDB,
		conversationClient: conversationClient,
		toRedisConsumer:    toRedisConsumer,
		toMongoConsumer:    toMongoConsumer,
		toMongoProducer:    toMongoProducer,
		toPushProducer:     toPushProducer,
	}, nil
}

// Start 启动消息传输处理
func (h *MsgTransferHandler) Start(ctx context.Context) error {
	// 启动在线消息处理（消费 toRedis topic → 分发到 Redis/Mongo/Push）
	go h.handleOnlineMsg(ctx)

	// 启动 toMongo 消费者
	go h.handleToMongo(ctx)

	logx.Info("MsgTransfer handler started")
	return nil
}

// handleOnlineMsg 处理在线消息（消费 toRedis topic → 分发到 Redis 缓存、MongoDB 持久化、Push 推送）
func (h *MsgTransferHandler) handleOnlineMsg(ctx context.Context) {
	onlineMsgHandler := NewOnlineMsgHandler(h.msgCache, h.conversationClient, h.toMongoProducer, h.toPushProducer)

	// 持续循环消费消息
	for {
		err := h.toRedisConsumer.Subscribe(ctx, func(msg mq.Message) error {
			return onlineMsgHandler.HandleMessage(msg)
		})
		if err != nil {
			logx.Errorf("Failed to subscribe to toRedis: %v", err)
			// 如果是 context 取消，退出循环
			select {
			case <-ctx.Done():
				return
			default:
				// 其他错误，继续尝试
				continue
			}
		}
	}
}

// handleToMongo 处理发送到 MongoDB 的消息
func (h *MsgTransferHandler) handleToMongo(ctx context.Context) {
	// 创建 ToMongoHandler
	toMongoHandler := NewToMongoHandler(h.msgDB)

	// 持续循环消费消息
	for {
		err := h.toMongoConsumer.Subscribe(ctx, func(msg mq.Message) error {
			// 使用 ToMongoHandler 处理消息
			return toMongoHandler.HandleMessage(msg)
		})
		if err != nil {
			logx.Errorf("Failed to subscribe to toMongo: %v", err)
			// 如果是 context 取消，退出循环
			select {
			case <-ctx.Done():
				return
			default:
				// 其他错误，继续尝试
				continue
			}
		}
	}
}

// Stop 停止消息传输处理
func (h *MsgTransferHandler) Stop(ctx context.Context) error {
	logx.Info("Stopping MsgTransfer handler...")

	var firstErr error

	// 1. 关闭 NATS 消费者（停止接收新消息）
	if h.toRedisConsumer != nil {
		if err := h.toRedisConsumer.Close(); err != nil {
			logx.Errorf("Failed to close toRedis consumer: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if h.toMongoConsumer != nil {
		if err := h.toMongoConsumer.Close(); err != nil {
			logx.Errorf("Failed to close toMongo consumer: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// 2. 关闭 NATS 生产者
	if h.toMongoProducer != nil {
		if err := h.toMongoProducer.Close(); err != nil {
			logx.Errorf("Failed to close toMongo producer: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if h.toPushProducer != nil {
		if err := h.toPushProducer.Close(); err != nil {
			logx.Errorf("Failed to close toPush producer: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// 3. 关闭 Redis 连接
	if h.redisClient != nil {
		if err := h.redisClient.Close(); err != nil {
			logx.Errorf("Failed to close Redis client: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// 4. 关闭 MongoDB 连接
	if h.mongoDB != nil {
		if err := h.mongoDB.Close(ctx); err != nil {
			logx.Errorf("Failed to close MongoDB: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	logx.Info("MsgTransfer handler stopped")
	return firstErr
}
