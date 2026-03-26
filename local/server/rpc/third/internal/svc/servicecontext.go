package svc

import (
	"context"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/pkg/s3"
	"fdim/pkg/s3/minio"
	"fdim/rpc/third/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config        config.Config
	MongoDB       *database.MongoDB
	Redis         *redis.Client
	ObjectStorage s3.Interface
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Host, c.Mongo.Database, c.Mongo.Username, c.Mongo.Password)

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化 MinIO 对象存储
	var objectStorage s3.Interface
	if c.ObjectStorage.Type == "minio" {
		minioCache := minio.NewRedisCache(redisClient)
		minioConfig := minio.Config{
			Bucket:          c.ObjectStorage.Bucket,
			Endpoint:        c.ObjectStorage.Endpoint,
			SignEndpoint:    c.ObjectStorage.SignEndpoint,
			AccessKeyID:     c.ObjectStorage.AccessKeyID,
			SecretAccessKey: c.ObjectStorage.SecretAccessKey,
			PublicRead:      true,
		}
		m, err := minio.NewMinio(context.Background(), minioCache, minioConfig)
		if err != nil {
			logx.Errorf("Failed to initialize MinIO: %v", err)
			// 不赋值给 objectStorage，保持 nil interface，避免 nil pointer panic
		} else {
			objectStorage = m
			logx.Info("MinIO object storage initialized successfully")
		}
	}

	return &ServiceContext{
		Config:        c,
		MongoDB:       mongoDB,
		Redis:         redisClient,
		ObjectStorage: objectStorage,
	}
}
