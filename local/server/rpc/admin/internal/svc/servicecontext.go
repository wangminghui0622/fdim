package svc

import (
	"time"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/pkg/tokenverify"
	"fdim/rpc/admin/internal/config"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config  config.Config
	MongoDB *database.MongoDB
	Redis   *redis.Client
	AdminDB database.AdminDatabase
	Token   *tokenverify.Token
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Url, c.Mongo.Db, "", "")

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化 AdminDatabase
	adminDB := database.NewAdminDatabase(mongoDB, redisClient)

	// 初始化 Token
	tokenExpire := time.Duration(c.TokenPolicy.Expire) * 24 * time.Hour
	token := &tokenverify.Token{
		Expires: tokenExpire,
		Secret:  c.Secret,
	}

	return &ServiceContext{
		Config:  c,
		MongoDB: mongoDB,
		Redis:   redisClient,
		AdminDB: adminDB,
		Token:   token,
	}
}
