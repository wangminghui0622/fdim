package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GroupCache struct {
	client *redis.Client
}

func NewGroupCache(client *redis.Client) *GroupCache {
	return &GroupCache{client: client}
}

// CacheGroupInfo 缓存群组信息
func (c *GroupCache) CacheGroupInfo(ctx context.Context, groupID string, groupInfo interface{}) error {
	// 实现逻辑
	return nil
}

// GetGroupInfo 获取群组信息
func (c *GroupCache) GetGroupInfo(ctx context.Context, groupID string) (interface{}, error) {
	// 实现逻辑
	return nil, nil
}
