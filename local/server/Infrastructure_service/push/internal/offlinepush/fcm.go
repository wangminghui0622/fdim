package offlinepush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// FcmPusher 使用 FCM 进行离线推?
type FcmPusher struct {
	redis     *redis.Client
	serverKey string
	client    *http.Client
}

// NewFcmPusher 创建 FCM 推送器
func NewFcmPusher(redis *redis.Client, serverKey string) *FcmPusher {
	return &FcmPusher{
		redis:     redis,
		serverKey: serverKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// fcmLegacyEndpoint 使用 FCM legacy HTTP 接口
const fcmLegacyEndpoint = "https://fcm.googleapis.com/fcm/send"

type fcmNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type fcmData struct {
	Ex          string `json:"ex,omitempty"`
	ClientMsgID string `json:"clientMsgID,omitempty"`
}

type fcmRequest struct {
	To           string          `json:"to,omitempty"`
	Registration []string        `json:"registration_ids,omitempty"`
	Notification *fcmNotification `json:"notification,omitempty"`
	Data         *fcmData        `json:"data,omitempty"`
}

// Push 根据 userIDs 查找各自?FCM token，并发送通知
func (p *FcmPusher) Push(ctx context.Context, userIDs []string, title, content string, opts *Options) error {
	if len(userIDs) == 0 {
		return nil
	}
	if p.redis == nil || p.serverKey == "" {
		logx.Error("FcmPusher not properly initialized (redis or serverKey is nil)")
		return nil
	}

	// 1. ?Redis 读取每个用户?FCM token（支持多终端?
	tokens := make([]string, 0)
	for _, userID := range userIDs {
		prefix := fmt.Sprintf("fcm_token:%s:", userID)
		iter := p.redis.Scan(ctx, 0, prefix+"*", 100).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			val, err := p.redis.Get(ctx, key).Result()
			if err != nil {
				logx.Errorf("FcmPusher: failed to get token from redis, key=%s, err=%v", key, err)
				continue
			}
			val = strings.TrimSpace(val)
			if val != "" {
				tokens = append(tokens, val)
			}
		}
		if err := iter.Err(); err != nil {
			logx.Errorf("FcmPusher: redis scan error for prefix=%s, err=%v", prefix, err)
		}
	}

	if len(tokens) == 0 {
		logx.Infof("FcmPusher: no tokens found for userIDs=%v", userIDs)
		return nil
	}

	// 2. 构?FCM 请求?
	data := &fcmData{}
	if opts != nil {
		data.Ex = opts.Ex
		if opts.Signal != nil {
			data.ClientMsgID = opts.Signal.ClientMsgID
		}
	}

	reqBody := fcmRequest{
		Registration: tokens,
		Notification: &fcmNotification{
			Title: title,
			Body:  content,
		},
		Data: data,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logx.Errorf("FcmPusher: failed to marshal request body: %v", err)
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fcmLegacyEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		logx.Errorf("FcmPusher: failed to create http request: %v", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "key="+p.serverKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		logx.Errorf("FcmPusher: http request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logx.Infof("FcmPusher: push success, tokens=%d", len(tokens))
		return nil
	}

	logx.Errorf("FcmPusher: push failed, status=%d", resp.StatusCode)
	return fmt.Errorf("fcm push failed, status=%d", resp.StatusCode)
}

