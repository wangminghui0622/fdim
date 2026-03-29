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
	ObjectStorage struct {
		Type           string // s3, minio, oss, etc.
		Endpoint       string // 内部连接地址（如 http://127.0.0.1:9000?
		SignEndpoint   string // 外部签名地址（如 http://42.192.129.44:9000），为空则与 Endpoint 相同
		Bucket         string
		AccessKeyID    string
		SecretAccessKey string
	}
}
