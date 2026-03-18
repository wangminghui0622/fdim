package svc

import (
	"time"

	"fdim/pkg/cache"
	"fdim/pkg/email"
	"fdim/pkg/sms"
	"fdim/pkg/database"
	"fdim/pkg/tokenverify"
	"fdim/rpc/chat/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"

	adminclient "fdim/protocol/admin"
)

type ServiceContext struct {
	Config   config.Config
	MongoDB  *database.MongoDB
	Redis    *redis.Client
	ChatDB   database.ChatDatabase
	Token    *tokenverify.Token
	Mailer   email.Mail
	SMS      sms.SMS
	AdminRpc adminclient.AdminClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MongoDB
	mongoDB := database.NewMongoDB(c.Mongo.Url, c.Mongo.Db, "", "")

	// 初始化 Redis
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始化 ChatDatabase
	chatDB := database.NewChatDatabase(mongoDB)

	// 初始化 Token
	tokenExpire := time.Duration(c.TokenPolicy.Expire) * 24 * time.Hour
	if tokenExpire == 0 {
		tokenExpire = 7 * 24 * time.Hour // 默认7天
	}
	token := &tokenverify.Token{
		Expires: tokenExpire,
		Secret:  c.Secret,
	}

	// 初始化 Admin RPC 客户端
	var adminRpc adminclient.AdminClient
	if c.AdminRpc.Etcd.Hosts != nil && len(c.AdminRpc.Etcd.Hosts) > 0 {
		adminRpc = adminclient.NewAdminClient(zrpc.MustNewClient(c.AdminRpc).Conn())
	}

	// 初始化邮件/短信发送器（可选）
	var mailer email.Mail
	if c.Email.Enable {
		mailer = email.NewMail(c.Email.SMTPAddr, c.Email.SMTPPort, c.Email.SenderMail, c.Email.SenderAuthorizationCode, c.Email.Title)
	}
	var smsSender sms.SMS
	if c.SMS.Enable {
		switch c.SMS.Provider {
		case "", "ali":
			aliSMS, err := sms.NewAli(c.SMS.Endpoint, c.SMS.AccessKeyID, c.SMS.AccessKeySecret, c.SMS.SignName, c.SMS.VerificationCodeTemplateCode)
			if err == nil {
				smsSender = aliSMS
			}
		}
	}

	return &ServiceContext{
		Config:   c,
		MongoDB:  mongoDB,
		Redis:    redisClient,
		ChatDB:   chatDB,
		Token:    token,
		Mailer:   mailer,
		SMS:      smsSender,
		AdminRpc: adminRpc,
	}
}
