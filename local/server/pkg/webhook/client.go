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

// Client Webhook ͻ??
type Client struct {
	url    string
	client *http.Client
}

// NewWebhookClient  Webhook ͻ??
func NewWebhookClient(url string) *Client {
	if url == "" {
		return nil //  URL Ϊգ??nilʾ webhook
	}
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallbackReq Webhook صӿ
type CallbackReq interface {
	GetCallbackCommand() string
}

// CallbackResp Webhook صӦӿ
type CallbackResp interface {
	Parse() error
}

// CommonCallbackResp ͨûصӦ
type CommonCallbackResp struct {
	ActionCode int32  `json:"actionCode"`
	ErrCode    int32  `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
	ErrDlt     string `json:"errDlt"`
	NextCode   int32  `json:"nextCode"`
}

const (
	// Next ִ
	Next = 1
	// NoError ޴??
	NoError = 0
	// CallbackError ص
	CallbackError = 1000
)

// ErrCallbackContinue ʾ webhook صʧܵӦüִ??
var ErrCallbackContinue = fmt.Errorf("webhook callback error but continue")

// Parse Ӧ
func (c *CommonCallbackResp) Parse() error {
	if c.ActionCode == NoError && c.NextCode == Next {
		//  ErrCode ??CallbackError??ErrCallbackContinueִ??
		if c.ErrCode == CallbackError {
			return ErrCallbackContinue
		}
		return fmt.Errorf("webhook callback error: code=%d, msg=%s, detail=%s", c.ErrCode, c.ErrMsg, c.ErrDlt)
	}
	return nil
}

// SyncPost ͬ??Webhook صBefore ص??
func (c *Client) SyncPost(ctx context.Context, command string, req CallbackReq, resp CallbackResp, timeout int) error {
	if c == nil {
		return nil // δ??webhook
	}

	fullURL := c.url + "/" + command
	operationID := mcontext.GetOperationID(ctx)

	logx.WithContext(ctx).Infof("webhook sync post: url=%s, command=%s, operationID=%s", fullURL, command, operationID)

	// ??
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook request: %w", err)
	}

	// 
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if operationID != "" {
		httpReq.Header.Set("operationID", operationID)
	}

	// óʱ
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
		httpReq = httpReq.WithContext(ctx)
	}

	// ??
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Ӧ
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("failed to decode webhook response: %w", err)
	}

	// 
	if err := resp.Parse(); err != nil {
		return err
	}

	logx.WithContext(ctx).Infof("webhook sync post success: url=%s", fullURL)
	return nil
}

// AsyncPost 첽??Webhook صAfter ص??
func (c *Client) AsyncPost(ctx context.Context, command string, req CallbackReq, resp CallbackResp, timeout int) {
	if c == nil {
		return // δ??webhook
	}

	// 첽ִ
	go func() {
		// ʹµ contextԭ context ȡ??
		asyncCtx := context.Background()
		if err := c.SyncPost(asyncCtx, command, req, resp, timeout); err != nil {
			logx.WithContext(asyncCtx).Errorf("webhook async post failed: %v", err)
		}
	}()
}

// WithCondition ִ Webhook ص
func WithCondition(ctx context.Context, enable bool, callback func(context.Context) error) error {
	if !enable {
		return nil
	}
	return callback(ctx)
}
