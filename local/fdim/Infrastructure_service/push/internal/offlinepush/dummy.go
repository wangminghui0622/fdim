package offlinepush

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

// DummyPusher 虚拟推送器（用于开发和测试）
type DummyPusher struct{}

// NewDummyPusher 创建虚拟推送器
func NewDummyPusher() *DummyPusher {
	return &DummyPusher{}
}

// Push 推送消息（虚拟实现，只记录日志）
func (d *DummyPusher) Push(ctx context.Context, userIDs []string, title, content string, opts *Options) error {
	logx.Infof("Dummy offline push: userIDs=%v, title=%s, content=%s", userIDs, title, content)
	// 虚拟推送器不做任何实际操作，只记录日志
	return nil
}
