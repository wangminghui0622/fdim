package config

import (
	"fdim/pkg/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo struct {
		Host     string
		Database string
		Username string
		Password string
	}
	Cache cache.CacheConf `yaml:"Cache"`
	Nats  struct {
		Brokers []string
		Topics  struct {
			ToRedis       string
			ToMongo       string
			ToPush        string
			ToOfflinePush string
		}
	}
	// RPC 客户端配
	ConversationRpc zrpc.RpcClientConf `yaml:"ConversationRpc" json:",optional"`
	// 通知配置（与官方一致）
	NotificationConfig config.Notification `yaml:"Notification" json:",optional"`
}
