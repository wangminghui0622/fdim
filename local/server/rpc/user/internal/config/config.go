package config

import (
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
	Cache cache.CacheConf `yaml:"Cache"` // 使用 Cache 作为字段名，配置文件使用 Cache ?
	// AdminUserIDs 管理员用户ID列表
	AdminUserIDs []string `json:",optional"`
	// MsgRpc msg RPC 服务配置，用于发送通知
	MsgRpc zrpc.RpcClientConf `yaml:"MsgRpc" json:",optional"`
	// ConversationRpc conversation RPC 服务配置，用于自动创建会?
	ConversationRpc zrpc.RpcClientConf `yaml:"ConversationRpc" json:",optional"`
	// Webhook Webhook 配置
	Webhook struct {
		URL string `json:",optional"` // Webhook 服务器地址
	}
}
