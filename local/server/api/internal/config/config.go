package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// RPC 服务配置
	UserRpc           zrpc.RpcClientConf
	FriendRpc         zrpc.RpcClientConf
	GroupRpc          zrpc.RpcClientConf
	AuthRpc           zrpc.RpcClientConf
	ConversationRpc   zrpc.RpcClientConf
	MsgRpc            zrpc.RpcClientConf
	ThirdRpc          zrpc.RpcClientConf
	MessageGatewayRpc zrpc.RpcClientConf `json:",optional"`
	ChatRpc           zrpc.RpcClientConf `json:",optional"`
	AdminRpc          zrpc.RpcClientConf `json:",optional"`
	// Secret 用于Token验证
	Secret string
	// MongoDB 配置
	Mongo struct {
		URI      string
		Database string
	} `json:",optional"`
	// AdminUserIDs 管理员用户ID列表
	AdminUserIDs []string `json:",optional"`
	// LiveKit 音视频通话配置
	LiveKit struct {
		URL    string `json:",optional"` // ws://your-server-ip:7880
		Key    string `json:",optional"` // API Key
		Secret string `json:",optional"` // API Secret
	} `json:",optional"`
}
