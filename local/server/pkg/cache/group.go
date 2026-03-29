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

// CacheGroupInfo ȺϢ
func (c *GroupCache) CacheGroupInfo(ctx context.Context, groupID string, groupInfo interface{}) error {
	// ʵ߼
	return nil
}

// GetGroupInfo ȡȺϢ
func (c *GroupCache) GetGroupInfo(ctx context.Context, groupID string) (interface{}, error) {
	// ʵ߼
	return nil, nil
}
