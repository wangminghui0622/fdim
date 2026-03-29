package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Name string `yaml:"Name"`
	Mode string `yaml:"Mode,default=dev"`
	Log  logx.LogConf `yaml:"Log"`

	// Cron 任务配置
	CronTask struct {
		// CronExecuteTime cron 执行时间表达式，例如 "0 0 2 * * *" 表示每天凌晨2点执?
		CronExecuteTime string `yaml:"CronExecuteTime"`
		// RetainChatRecords 保留聊天记录天数，超过此天数的记录将被删?
		RetainChatRecords int `yaml:"RetainChatRecords"`
		// FileExpireTime 文件过期时间（天），超过此时间的文件将被删除
		FileExpireTime int `yaml:"FileExpireTime"`
		// DeleteObjectType 要删除的对象类型列表
		DeleteObjectType []string `yaml:"DeleteObjectType"`
	} `yaml:"CronTask"`
	
	// RPC 客户端配?
	MsgRpc          zrpc.RpcClientConf `yaml:"MsgRpc"`
	ConversationRpc zrpc.RpcClientConf `yaml:"ConversationRpc"`
	ThirdRpc        zrpc.RpcClientConf `yaml:"ThirdRpc"`
	
	// 管理员用户ID（用于执行定时任务）
	AdminUserIDs []string `yaml:"AdminUserIDs"`
}
