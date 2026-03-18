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

// CacheFriendList 缓存好友列表
func (c *FriendCache) CacheFriendList(ctx context.Context, ownerUserID string, friendIDs []string) error {
	// 实现逻辑
	return nil
}

// GetFriendList 获取好友列表
func (c *FriendCache) GetFriendList(ctx context.Context, ownerUserID string) ([]string, error) {
	// 实现逻辑
	return []string{}, nil
}
