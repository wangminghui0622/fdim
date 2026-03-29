package config

import "github.com/zeromicro/go-zero/zrpc"

type RedisConfig struct {
	Host string `yaml:"Host"`
	Type string `yaml:"Type"`
	Pass string `yaml:"Pass"`
}

type Config struct {
	Name string `yaml:"Name"`
	Mode string `yaml:"Mode"`
	// MongoDB 配置
	MongoDB struct {
		URI      string `yaml:"URI"`
		Database string `yaml:"Database"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"MongoDB"`
	
	// Redis 配置
	Cache []RedisConfig `yaml:"Cache"`
	
	// NATS 配置
	Nats struct {
		Brokers            []string `yaml:"Brokers"`
		ToRedisTopic       string   `yaml:"ToRedisTopic"`
		ToMongoTopic      string   `yaml:"ToMongoTopic"`
		ToPushTopic       string   `yaml:"ToPushTopic"`
		ToOfflinePushTopic string  `yaml:"ToOfflinePushTopic"`
		ToRedisGroupID    string   `yaml:"ToRedisGroupID"`
		ToMongoGroupID    string   `yaml:"ToMongoGroupID"`
		ToPushGroupID     string   `yaml:"ToPushGroupID"`
		ToOfflinePushGroupID string `yaml:"ToOfflinePushGroupID"`
	} `yaml:"Nats"`
	
	// RPC 客户端配?
	ConversationRpc zrpc.RpcClientConf `yaml:"ConversationRpc" json:",optional"`
	
	// 批处理配?
	Batch struct {
		Size     int `yaml:"Size"`     // 批处理大?
		Interval int `yaml:"Interval"` // 批处理间隔（毫秒?
		Workers  int `yaml:"Workers"`  // 工作协程?
	} `yaml:"Batch"`
}
