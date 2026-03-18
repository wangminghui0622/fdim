package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fdim/pkg/mcontext"
	"github.com/zeromicro/go-zero/core/logx"
)

// Client Webhook 客户�?
type Client struct {
	url    string
	client *http.Client
}

// NewWebhookClient 创建 Webhook 客户�?
func NewWebhookClient(url string) *Client {
	if url == "" {
		return nil // 如果 URL 为空，返�?nil，表示不启用 webhook
	}
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallbackReq Webhook 回调请求接口
type CallbackReq interface {
	GetCallbackCommand() string
}

// CallbackResp Webhook 回调响应接口
type CallbackResp interface {
	Parse() error
}

// CommonCallbackResp 通用回调响应
type CommonCallbackResp struct {
	ActionCode int32  `json:"actionCode"`
	ErrCode    int32  `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
	ErrDlt     string `json:"errDlt"`
	NextCode   int32  `json:"nextCode"`
}

const (
	// Next 继续执行
	Next = 1
	// NoError 无错�?
	NoError = 0
	// CallbackError 回调错误代码
	CallbackError = 1000
)

// ErrCallbackContinue 表示 webhook 回调失败但应该继续执�?
var ErrCallbackContinue = fmt.Errorf("webhook callback error but continue")

// Parse 解析响应
func (c *CommonCallbackResp) Parse() error {
	if c.ActionCode == NoError && c.NextCode == Next {
		// 如果 ErrCode �?CallbackError，返�?ErrCallbackContinue，允许继续执�?
		if c.ErrCode == CallbackError {
			return ErrCallbackContinue
		}
		return fmt.Errorf("webhook callback error: code=%d, msg=%s, detail=%s", c.ErrCode, c.ErrMsg, c.ErrDlt)
	}
	return nil
}

// SyncPost 同步发�?Webhook 回调（Before 回调�?
func (c *Client) SyncPost(ctx context.Context, command string, req CallbackReq, resp CallbackResp, timeout int) error {
	if c == nil {
		return nil // 未启�?webhook
	}

	fullURL := c.url + "/" + command
	operationID := mcontext.GetOperationID(ctx)

	logx.WithContext(ctx).Infof("webhook sync post: url=%s, command=%s, operationID=%s", fullURL, command, operationID)

	// 创建请求�?
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook request: %w", err)
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if operationID != "" {
		httpReq.Header.Set("operationID", operationID)
	}

	// 设置超时
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}

	// 发送请�?
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// 解析响应
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("failed to decode webhook response: %w", err)
	}

	// 解析错误
	if err := resp.Parse(); err != nil {
		return err
	}

	logx.WithContext(ctx).Infof("webhook sync post success: url=%s", fullURL)
	return nil
}

// AsyncPost 异步发�?Webhook 回调（After 回调�?
func (c *Client) AsyncPost(ctx context.Context, command string, req CallbackReq, resp CallbackResp, timeout int) {
	if c == nil {
		return // 未启�?webhook
	}

	// 异步执行
	go func() {
		// 使用新的 context，避免原 context 被取�?
		asyncCtx := context.Background()
		if err := c.SyncPost(asyncCtx, command, req, resp, timeout); err != nil {
			logx.WithContext(asyncCtx).Errorf("webhook async post failed: %v", err)
		}
	}()
}

// WithCondition 条件执行 Webhook 回调
func WithCondition(ctx context.Context, enable bool, callback func(context.Context) error) error {
	if !enable {
		return nil
	}
	return callback(ctx)
}
