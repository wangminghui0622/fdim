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
	Cache cache.CacheConf `yaml:"Cache"`
	// Secret 用于Token签名
	Secret string
	// TokenPolicy Token策略
	TokenPolicy struct {
		Expire int // 过期天数
	}
	// MultiLogin 是否允许多设备登录
	MultiLogin bool
	// AdminUserIDs 管理员用户ID列表
	AdminUserIDs []string `json:",optional"`

	// RPC 客户端配置
	Rpc struct {
		MsgGateway zrpc.RpcClientConf `yaml:"MsgGateway"`
	} `yaml:"Rpc"`
}
