package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo struct {
		Url string
		Db  string
	}
	Cache cache.CacheConf `yaml:"Cache"`

	// RPC 客户端配置
	MsgRpc  zrpc.RpcClientConf `yaml:"MsgRpc" json:",optional"`
	UserRpc zrpc.RpcClientConf `yaml:"UserRpc" json:",optional"`
}
