package config

import "strings"

var EnvPrefixMap map[string]string

func init() {
	EnvPrefixMap = make(map[string]string)
	fileNames := []string{
		FileName, NotificationFileName, ShareFileName, WebhooksConfigFileName,
		NatsConfigFileName, RedisConfigFileName,
		MongodbConfigFileName, MinioConfigFileName, LogConfigFileName,
		FDIMAPICfgFileName, FDIMCronTaskCfgFileName, FDIMMsgGatewayCfgFileName,
		FDIMMsgTransferCfgFileName, FDIMPushCfgFileName, FDIMRPCAuthCfgFileName,
		FDIMRPCConversationCfgFileName, FDIMRPCFriendCfgFileName, FDIMRPCGroupCfgFileName,
		FDIMRPCMsgCfgFileName, FDIMRPCThirdCfgFileName, FDIMRPCUserCfgFileName, DiscoveryConfigFilename,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "IMENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		EnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
