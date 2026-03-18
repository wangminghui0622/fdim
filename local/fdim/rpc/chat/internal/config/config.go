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
	AdminRpc zrpc.RpcClientConf `yaml:"AdminRpc"`
	// Secret 用于Token签名
	Secret string `yaml:"Secret"`
	// TokenPolicy Token策略
	TokenPolicy struct {
		Expire int `yaml:"Expire"` // 天数
	} `yaml:"TokenPolicy"`
	// AllowRegister 是否允许注册
	AllowRegister bool `yaml:"AllowRegister,default=true"`
	Email         struct {
		Enable                    bool   `yaml:"Enable"`
		SMTPAddr                  string `yaml:"SMTPAddr"`
		SMTPPort                  int    `yaml:"SMTPPort"`
		SenderMail                string `yaml:"SenderMail"`
		SenderAuthorizationCode   string `yaml:"SenderAuthorizationCode"`
		Title                     string `yaml:"Title"`
	} `yaml:"Email"`
	SMS struct {
		Enable                         bool   `yaml:"Enable"`
		Provider                       string `yaml:"Provider"`
		Endpoint                       string `yaml:"Endpoint"`
		AccessKeyID                    string `yaml:"AccessKeyID"`
		AccessKeySecret                string `yaml:"AccessKeySecret"`
		SignName                       string `yaml:"SignName"`
		VerificationCodeTemplateCode   string `yaml:"VerificationCodeTemplateCode"`
	} `yaml:"SMS"`
	RTC struct {
		ServerURL string `yaml:"ServerURL"`
		TokenSalt string `yaml:"TokenSalt"`
	} `yaml:"RTC"`
}
