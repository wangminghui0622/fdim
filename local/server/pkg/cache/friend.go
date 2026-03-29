package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type FriendCache struct {
	client *redis.Client
}

func NewFriendCache(client *redis.Client) *FriendCache {
	return &FriendCache{client: client}
}

// CacheFriendList б
func (c *FriendCache) CacheFriendList(ctx context.Context, ownerUserID string, friendIDs []string) error {
	// ʵ߼
	return nil
}

// GetFriendList ȡб
func (c *FriendCache) GetFriendList(ctx context.Context, ownerUserID string) ([]string, error) {
	// ʵ߼
	return []string{}, nil
}
