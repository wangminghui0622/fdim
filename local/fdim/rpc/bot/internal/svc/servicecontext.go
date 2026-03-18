package svc

import (
	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/protocol/msg"
	"fdim/rpc/bot/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	MongoDB *database.MongoDB
	Redis   *redis.Client
	BotDB   database.BotDatabase
	MsgRpc  msg.MsgClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Url, c.Mongo.Db, "", "")

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化 BotDatabase
	botDB := database.NewBotDatabase(mongoDB)

	// 初始化 Msg RPC 客户端
	var msgRpc msg.MsgClient
	if c.MsgRpc.Etcd.Hosts != nil && len(c.MsgRpc.Etcd.Hosts) > 0 {
		msgRpc = msg.NewMsgClient(zrpc.MustNewClient(c.MsgRpc).Conn())
	}

	return &ServiceContext{
		Config:  c,
		MongoDB: mongoDB,
		Redis:   redisClient,
		BotDB:   botDB,
		MsgRpc:  msgRpc,
	}
}
