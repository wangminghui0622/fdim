package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fdim/pkg/model"
	"github.com/redis/go-redis/v9"
)

// MsgCache 消息缓存接口
type MsgCache interface {
	// IncrMaxSeq 原子递增会话序列号，返回新的 seq（不写消息缓存）
	IncrMaxSeq(ctx context.Context, conversationID string) (int64, error)
	// SetMessagesToCache 将消息写入缓存（使用消息自带的 seq，不分配新 seq）
	SetMessagesToCache(ctx context.Context, conversationID string, msgs []*model.MsgDoc) error
	// BatchInsertChat2Cache 批量插入聊天消息到缓存（分配新 seq + 写入）
	BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*model.MsgDoc) (int64, bool, map[string]int64, error)
	// GetMessagesBySeq 根据序列号获取消息
	GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgDoc, error)
	// GetMaxSeq 获取最大序列号
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	// SetMaxSeq 设置最大序列号（用于对齐会话序列号）
	SetMaxSeq(ctx context.Context, conversationID string, seq int64) error
	// GetMinSeq 获取最小保留序列号
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	// SetMinSeq 设置最小保留序列号（用于清理历史消息）
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error
	// SetHasReadSeqs 设置已读序列号
	SetHasReadSeqs(ctx context.Context, conversationID string, userSeqMap map[string]int64) error
	// GetHasReadSeq 获取单个用户在会话中的已读序列号
	GetHasReadSeq(ctx context.Context, conversationID, userID string) (int64, error)
}

// RedisMsgCache Redis消息缓存实现
type RedisMsgCache struct {
	client *redis.Client
	msgDB  interface{} // TODO: 添加消息数据库
}

// NewRedisMsgCache 创建Redis消息缓存
func NewRedisMsgCache(client *redis.Client, msgDB interface{}) *RedisMsgCache {
	return &RedisMsgCache{
		client: client,
		msgDB:  msgDB,
	}
}

// getConversationMaxSeqKey 获取会话最大序列号键
func (c *RedisMsgCache) getConversationMaxSeqKey(conversationID string) string {
	return fmt.Sprintf("conversation:max_seq:%s", conversationID)
}

// getConversationMinSeqKey 获取会话最小序列号键
func (c *RedisMsgCache) getConversationMinSeqKey(conversationID string) string {
	return fmt.Sprintf("conversation:min_seq:%s", conversationID)
}

// getMessageKey 获取消息键
func (c *RedisMsgCache) getMessageKey(conversationID string, seq int64) string {
	return fmt.Sprintf("msg:%s:%d", conversationID, seq)
}

// IncrMaxSeq 原子递增会话序列号，返回新的 seq（不写消息缓存）
func (c *RedisMsgCache) IncrMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	maxSeqKey := c.getConversationMaxSeqKey(conversationID)
	newSeq, err := c.client.Incr(ctx, maxSeqKey).Result()
	if err != nil {
		return 0, fmt.Errorf("递增序列号失败: %w", err)
	}
	return newSeq, nil
}

// SetMessagesToCache 将消息写入缓存（使用消息自带的 seq，不分配新 seq）
func (c *RedisMsgCache) SetMessagesToCache(ctx context.Context, conversationID string, msgs []*model.MsgDoc) error {
	if len(msgs) == 0 {
		return nil
	}
	pipe := c.client.Pipeline()
	for _, msg := range msgs {
		msgKey := c.getMessageKey(conversationID, msg.Seq)
		msgBytes, _ := json.Marshal(msg)
		pipe.Set(ctx, msgKey, msgBytes, 7*24*time.Hour)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("写入消息缓存失败: %w", err)
	}
	return nil
}

// BatchInsertChat2Cache 批量插入聊天消息到缓存（分配新 seq + 写入）
func (c *RedisMsgCache) BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*model.MsgDoc) (int64, bool, map[string]int64, error) {
	if len(msgs) == 0 {
		return 0, false, nil, nil
	}

	pipe := c.client.Pipeline()

	// 获取当前最大序列号
	maxSeqKey := c.getConversationMaxSeqKey(conversationID)
	currentMaxSeq, err := c.client.Get(ctx, maxSeqKey).Int64()
	if err == redis.Nil {
		currentMaxSeq = 0
	} else if err != nil {
		return 0, false, nil, fmt.Errorf("获取最大序列号失败: %w", err)
	}

	isNewConversation := currentMaxSeq == 0

	// 存储消息并更新序列号
	userSeqMap := make(map[string]int64)
	lastSeq := currentMaxSeq

	for _, msg := range msgs {
		lastSeq++
		msg.Seq = lastSeq

		// 存储消息
		msgKey := c.getMessageKey(conversationID, lastSeq)
		msgBytes, _ := json.Marshal(msg)
		pipe.Set(ctx, msgKey, msgBytes, 7*24*time.Hour) // 7天过期

		// 更新用户序列号
		if msg.SendID != "" {
			userSeqMap[msg.SendID] = lastSeq
		}
	}

	// 更新最大序列号
	pipe.Set(ctx, maxSeqKey, lastSeq, 0) // 不过期

	// 执行管道操作
	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, false, nil, fmt.Errorf("执行管道操作失败: %w", err)
	}

	return lastSeq, isNewConversation, userSeqMap, nil
}

// GetMessagesBySeq 根据序列号获取消息
func (c *RedisMsgCache) GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgDoc, error) {
	if len(seqs) == 0 {
		return nil, nil
	}

	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(seqs))

	for i, seq := range seqs {
		msgKey := c.getMessageKey(conversationID, seq)
		cmds[i] = pipe.Get(ctx, msgKey)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("获取消息失败: %w", err)
	}

	var msgs []*model.MsgDoc
	for i, cmd := range cmds {
		val, err := cmd.Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("获取消息序列号%d失败: %w", seqs[i], err)
		}

		var msg model.MsgDoc
		if err := json.Unmarshal([]byte(val), &msg); err != nil {
			return nil, fmt.Errorf("反序列化消息失败: %w", err)
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

// GetMaxSeq 获取最大序列号
func (c *RedisMsgCache) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	maxSeqKey := c.getConversationMaxSeqKey(conversationID)
	seq, err := c.client.Get(ctx, maxSeqKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return seq, err
}

// SetMaxSeq 设置最大序列号
func (c *RedisMsgCache) SetMaxSeq(ctx context.Context, conversationID string, seq int64) error {
	maxSeqKey := c.getConversationMaxSeqKey(conversationID)
	return c.client.Set(ctx, maxSeqKey, seq, 0).Err()
}

// GetMinSeq 获取最小保留序列号
func (c *RedisMsgCache) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	minSeqKey := c.getConversationMinSeqKey(conversationID)
	seq, err := c.client.Get(ctx, minSeqKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return seq, err
}

// SetMinSeq 设置最小保留序列号
func (c *RedisMsgCache) SetMinSeq(ctx context.Context, conversationID string, seq int64) error {
	minSeqKey := c.getConversationMinSeqKey(conversationID)
	return c.client.Set(ctx, minSeqKey, seq, 0).Err()
}

// SetHasReadSeqs 设置已读序列号
func (c *RedisMsgCache) SetHasReadSeqs(ctx context.Context, conversationID string, userSeqMap map[string]int64) error {
	pipe := c.client.Pipeline()

	for userID, seq := range userSeqMap {
		key := fmt.Sprintf("conversation:has_read:%s:%s", conversationID, userID)
		pipe.Set(ctx, key, seq, 0)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetHasReadSeq 获取单个用户在会话中的已读序列号
func (c *RedisMsgCache) GetHasReadSeq(ctx context.Context, conversationID, userID string) (int64, error) {
	key := fmt.Sprintf("conversation:has_read:%s:%s", conversationID, userID)
	seq, err := c.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return seq, err
}
