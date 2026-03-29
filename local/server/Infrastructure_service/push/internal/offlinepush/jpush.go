package offlinepush

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// JPushPusher 使用极光推送进行离线推送
type JPushPusher struct {
	redis     *redis.Client
	appKey    string
	masterSecret string
	client    *http.Client
}

// NewJPushPusher 创建极光推送器
func NewJPushPusher(redis *redis.Client, appKey, masterSecret string) *JPushPusher {
	return &JPushPusher{
		redis:        redis,
		appKey:       appKey,
		masterSecret: masterSecret,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

const jpushEndpoint = "https://api.jpush.cn/v3/push"

type jpushPlatform []string

type jpushAudience struct {
	Alias        []string `json:"alias,omitempty"`
	RegistrationID []string `json:"registration_id,omitempty"`
}

type jpushNotification struct {
	Alert   string                 `json:"alert,omitempty"`
	Android *jpushAndroidNotification `json:"android,omitempty"`
	IOS     *jpushIOSNotification     `json:"ios,omitempty"`
}

type jpushAndroidNotification struct {
	Alert  string            `json:"alert,omitempty"`
	Title  string            `json:"title,omitempty"`
	Extras map[string]string `json:"extras,omitempty"`
}

type jpushIOSNotification struct {
	Alert  interface{}       `json:"alert,omitempty"`
	Sound  string            `json:"sound,omitempty"`
	Badge  interface{}       `json:"badge,omitempty"`
	Extras map[string]string `json:"extras,omitempty"`
}

type jpushMessage struct {
	MsgContent  string            `json:"msg_content,omitempty"`
	Title       string            `json:"title,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Extras      map[string]string `json:"extras,omitempty"`
}

type jpushOptions struct {
	TimeToLive     int  `json:"time_to_live,omitempty"`
	ApnsProduction bool `json:"apns_production,omitempty"`
}

type jpushRequest struct {
	Platform     interface{}        `json:"platform"`
	Audience     interface{}        `json:"audience"`
	Notification *jpushNotification `json:"notification,omitempty"`
	Message      *jpushMessage      `json:"message,omitempty"`
	Options      *jpushOptions      `json:"options,omitempty"`
}

// Push 根据 userIDs 查找各自的极光推?registration_id，并发送通知
func (p *JPushPusher) Push(ctx context.Context, userIDs []string, title, content string, opts *Options) error {
	if len(userIDs) == 0 {
		return nil
	}
	if p.redis == nil || p.appKey == "" || p.masterSecret == "" {
		logx.Error("JPushPusher not properly initialized")
		return nil
	}

	// 1. ?Redis 读取每个用户的极光推?registration_id
	registrationIDs := make([]string, 0)
	for _, userID := range userIDs {
		key := fmt.Sprintf("jpush_token:%s", userID)
		val, err := p.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		if val != "" {
			registrationIDs = append(registrationIDs, val)
		}
	}

	if len(registrationIDs) == 0 {
		logx.Infof("JPushPusher: no registration_ids found for userIDs=%v", userIDs)
		return nil
	}

	// 2. 构造极光推送请求体
	extras := make(map[string]string)
	if opts != nil {
		if opts.Ex != "" {
			extras["ex"] = opts.Ex
		}
		if opts.Signal != nil && opts.Signal.ClientMsgID != "" {
			extras["clientMsgID"] = opts.Signal.ClientMsgID
		}
	}

	reqBody := jpushRequest{
		Platform: "all",
		Audience: jpushAudience{
			RegistrationID: registrationIDs,
		},
		Notification: &jpushNotification{
			Alert: content,
			Android: &jpushAndroidNotification{
				Alert:  content,
				Title:  title,
				Extras: extras,
			},
			IOS: &jpushIOSNotification{
				Alert:  content,
				Sound:  "default",
				Badge:  "+1",
				Extras: extras,
			},
		},
		Options: &jpushOptions{
			TimeToLive:     86400, // 24 hours
			ApnsProduction: true,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logx.Errorf("JPushPusher: failed to marshal request body: %v", err)
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, jpushEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		logx.Errorf("JPushPusher: failed to create http request: %v", err)
		return err
	}

	// 设置 Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(p.appKey + ":" + p.masterSecret))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+auth)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		logx.Errorf("JPushPusher: http request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logx.Infof("JPushPusher: push success, registration_ids=%d", len(registrationIDs))
		return nil
	}

	logx.Errorf("JPushPusher: push failed, status=%d", resp.StatusCode)
	return fmt.Errorf("jpush push failed, status=%d", resp.StatusCode)
}
