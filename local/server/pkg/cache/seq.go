package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// SeqUser 用户序列号缓存接?
type SeqUser interface {
	// GetUserMaxSeq 获取用户最大序列号
	GetUserMaxSeq(ctx context.Context, userID string) (int64, error)
	// SetUserMaxSeq 设置用户最大序列号
	SetUserMaxSeq(ctx context.Context, userID string, seq int64) error
	// IncrUserSeq 递增用户序列?
	IncrUserSeq(ctx context.Context, userID string) (int64, error)
}

// SeqConversationCache 会话序列号缓存接?
type SeqConversationCache interface {
	// GetConversationMaxSeq 获取会话最大序列号
	GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error)
	// SetConversationMaxSeq 设置会话最大序列号
	SetConversationMaxSeq(ctx context.Context, conversationID string, seq int64) error
	// IncrConversationSeq 递增会话序列?
	IncrConversationSeq(ctx context.Context, conversationID string) (int64, error)
}

// RedisSeqUser Redis用户序列号缓存实?
type RedisSeqUser struct {
	client *redis.Client
}

// NewRedisSeqUser 创建Redis用户序列号缓?
func NewRedisSeqUser(client *redis.Client) *RedisSeqUser {
	return &RedisSeqUser{
		client: client,
	}
}

func (s *RedisSeqUser) getUserSeqKey(userID string) string {
	return fmt.Sprintf("user:seq:%s", userID)
}

// GetUserMaxSeq 获取用户最大序列号
func (s *RedisSeqUser) GetUserMaxSeq(ctx context.Context, userID string) (int64, error) {
	key := s.getUserSeqKey(userID)
	seq, err := s.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return seq, err
}

// SetUserMaxSeq 设置用户最大序列号
func (s *RedisSeqUser) SetUserMaxSeq(ctx context.Context, userID string, seq int64) error {
	key := s.getUserSeqKey(userID)
	return s.client.Set(ctx, key, seq, 0).Err()
}

// IncrUserSeq 增加用户序列?
func (s *RedisSeqUser) IncrUserSeq(ctx context.Context, userID string) (int64, error) {
	key := s.getUserSeqKey(userID)
	return s.client.Incr(ctx, key).Result()
}

// RedisSeqConversationCache Redis会话序列号缓存实?
type RedisSeqConversationCache struct {
	client *redis.Client
}

// NewRedisSeqConversationCache 创建Redis会话序列号缓?
func NewRedisSeqConversationCache(client *redis.Client) *RedisSeqConversationCache {
	return &RedisSeqConversationCache{
		client: client,
	}
}

func (s *RedisSeqConversationCache) getConversationSeqKey(conversationID string) string {
	return fmt.Sprintf("conversation:seq:%s", conversationID)
}

// GetConversationMaxSeq 获取会话最大序列号
func (s *RedisSeqConversationCache) GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	key := s.getConversationSeqKey(conversationID)
	seq, err := s.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return seq, err
}

// SetConversationMaxSeq 设置会话最大序列号
func (s *RedisSeqConversationCache) SetConversationMaxSeq(ctx context.Context, conversationID string, seq int64) error {
	key := s.getConversationSeqKey(conversationID)
	return s.client.Set(ctx, key, seq, 0).Err()
}

// IncrConversationSeq 增加会话序列?
func (s *RedisSeqConversationCache) IncrConversationSeq(ctx context.Context, conversationID string) (int64, error) {
	key := s.getConversationSeqKey(conversationID)
	return s.client.Incr(ctx, key).Result()
}
