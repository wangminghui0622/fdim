package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf `yaml:"Cache"`
	Mongo struct {
		Url string
		Db  string
	}
	Rpc struct {
		User zrpc.RpcClientConf `yaml:"User"`
		Auth zrpc.RpcClientConf `yaml:"Auth"`
	} `yaml:"Rpc"`
	// Secret 用于Token签名
	Secret string `json:",optional"`
	// TokenPolicy Token策略
	TokenPolicy struct {
		Expire int `json:",optional"` // 过期天数
	} `json:",optional"`
	// MultiLogin 是否允许多设备登录
	MultiLogin bool `json:",optional"`
	// AdminUserIDs 管理员用户ID列表
	AdminUserIDs []string `json:",optional"`
}
