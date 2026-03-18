package svc

import (
	"time"

	"fdim/pkg/cache"
	"fdim/pkg/database"
	"fdim/pkg/tokenverify"
	"fdim/protocol/msggateway"
	"fdim/rpc/auth/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	MongoDB          *database.MongoDB
	Redis            *redis.Client
	AuthDB           database.AuthDatabase
	Token            *tokenverify.Token
	MsgGatewayClient msggateway.MsgGatewayClient
	UserDB           *database.UserDatabase
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Host, c.Mongo.Database, c.Mongo.Username, c.Mongo.Password)

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化 Token 签名器
	tokenExpire := time.Duration(c.TokenPolicy.Expire) * 24 * time.Hour
	if tokenExpire == 0 {
		tokenExpire = 7 * 24 * time.Hour
	}
	token := &tokenverify.Token{
		Expires: tokenExpire,
		Secret:  c.Secret,
	}

	// 初始化 AuthDatabase
	authDB := database.NewAuthDatabase(redisClient, token, c.MultiLogin)

	// 初始化 MsgGateway 客户端（可选）
	var msgGatewayClient msggateway.MsgGatewayClient
	if c.Rpc.MsgGateway.Etcd.Key != "" || c.Rpc.MsgGateway.Target != "" {
		msgGatewayConn := zrpc.MustNewClient(c.Rpc.MsgGateway).Conn()
		msgGatewayClient = msggateway.NewMsgGatewayClient(msgGatewayConn)
	}

	// 初始化 UserDatabase（用于验证用户是否存在）
	userDB := database.NewUserDatabase(mongoDB)

	return &ServiceContext{
		Config:           c,
		MongoDB:          mongoDB,
		Redis:            redisClient,
		AuthDB:           authDB,
		Token:            token,
		MsgGatewayClient: msgGatewayClient,
		UserDB:           userDB,
	}
}
