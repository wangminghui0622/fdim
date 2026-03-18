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
	MsgRpc zrpc.RpcClientConf `yaml:"MsgRpc"`
}
