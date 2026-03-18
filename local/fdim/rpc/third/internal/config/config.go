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
		Type     string // s3, minio, oss, etc.
		Endpoint string
		Bucket   string
		AccessKeyID string
		SecretAccessKey string
	}
}
