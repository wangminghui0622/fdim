package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf `json:"Cache,optional"`
	
	// NATS 配置
	Nats struct {
		Brokers              []string `json:"Brokers"`
		ToPushTopic          string   `json:"ToPushTopic"`
		ToOfflinePushTopic   string   `json:"ToOfflinePushTopic"`
		ToPushGroupID        string   `json:"ToPushGroupID"`
		ToOfflinePushGroupID string   `json:"ToOfflinePushGroupID"`
	} `json:"Nats"`
	
	// RPC 客户端配?
	MessageGatewayRpc zrpc.RpcClientConf `json:"MessageGatewayRpc,optional"`
	GroupRpc          zrpc.RpcClientConf `json:"GroupRpc,optional"`
	
	// 离线推送配?
	Push struct {
		Enable            string `json:"Enable,optional"`            // getui, fcm, jpush, 或空（使?dummy?
		FcmServerKey      string `json:"FcmServerKey,optional"`      // FCM 服务器密?
		JPushAppKey       string `json:"JPushAppKey,optional"`       // 极光推?AppKey
		JPushMasterSecret string `json:"JPushMasterSecret,optional"` // 极光推?MasterSecret
	} `json:"Push,optional"`
}
