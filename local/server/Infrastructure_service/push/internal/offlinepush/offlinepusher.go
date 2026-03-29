package offlinepush

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// OfflinePusher 离线推送器接口
type OfflinePusher interface {
	Push(ctx context.Context, userIDs []string, title, content string, opts *Options) error
}

// Options 离线推送选项
type Options struct {
	Ex            string
	IOSPushSound  string
	IOSBadgeCount int32
	Signal        *Signal
}

// Signal 信号数据
type Signal struct {
	ClientMsgID string
}

// PushConfig 离线推送配?
type PushConfig struct {
	Enable            string
	FcmServerKey      string
	JPushAppKey       string
	JPushMasterSecret string
}

// NewOfflinePusher 创建离线推送器
// enable: getui, fcm, jpush, 或空（使?dummy?
// redis: 用于?Redis 中读取各用户的推?token
// fcmServerKey: FCM 的服务器密钥
func NewOfflinePusher(enable string, redis *redis.Client, fcmServerKey string) (OfflinePusher, error) {
	return NewOfflinePusherWithConfig(PushConfig{
		Enable:       enable,
		FcmServerKey: fcmServerKey,
	}, redis)
}

// NewOfflinePusherWithConfig 使用完整配置创建离线推送器
func NewOfflinePusherWithConfig(config PushConfig, redis *redis.Client) (OfflinePusher, error) {
	switch config.Enable {
	case "getui":
		// TODO: 实现 GetUI 推?
		return NewDummyPusher(), nil
	case "fcm":
		if redis == nil || config.FcmServerKey == "" {
			return NewDummyPusher(), nil
		}
		return NewFcmPusher(redis, config.FcmServerKey), nil
	case "jpush":
		if redis == nil || config.JPushAppKey == "" || config.JPushMasterSecret == "" {
			return NewDummyPusher(), nil
		}
		return NewJPushPusher(redis, config.JPushAppKey, config.JPushMasterSecret), nil
	default:
		return NewDummyPusher(), nil
	}
}
