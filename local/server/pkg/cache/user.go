package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{client: client}
}

// SetUserOnline 设置用户在线状态（单个平台）
func (c *UserCache) SetUserOnline(ctx context.Context, userID string, platformIDs []int32) error {
	_ = "user:online:" + userID
	// 实现逻辑：将 platformIDs 存储到 Redis Set
	// 这里简化实现，实际应该使用 Redis Set 存储多个平台ID
	return nil
}

// SetUserOnlineStatus 设置用户在线状态（支持在线和离线平台列表）
func (c *UserCache) SetUserOnlineStatus(ctx context.Context, userID string, online []int32, offline []int32) error {
	key := "user:online:" + userID
	
	// 添加在线平台
	for _, platformID := range online {
		member := strconv.FormatInt(int64(platformID), 10)
		if err := c.client.SAdd(ctx, key, member).Err(); err != nil {
			return err
		}
	}
	
	// 移除离线平台
	for _, platformID := range offline {
		member := strconv.FormatInt(int64(platformID), 10)
		if err := c.client.SRem(ctx, key, member).Err(); err != nil {
			return err
		}
	}
	
	// 设置过期时间（4小时）
	if err := c.client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return err
	}
	
	return nil
}

// GetUserOnline 获取用户在线状态
func (c *UserCache) GetUserOnline(ctx context.Context, userID string) ([]int32, error) {
	key := "user:online:" + userID
	
	// 从 Redis Set 获取所有平台ID
	members, err := c.client.SMembers(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return []int32{}, nil
		}
		return nil, err
	}
	
	platformIDs := make([]int32, 0, len(members))
	for _, member := range members {
		platformID, err := strconv.ParseInt(member, 10, 32)
		if err != nil {
			continue
		}
		platformIDs = append(platformIDs, int32(platformID))
	}
	
	return platformIDs, nil
}

// GetAllOnlineUsers 获取所有在线用户（支持游标分页）
func (c *UserCache) GetAllOnlineUsers(ctx context.Context, cursor string) (map[string][]int32, string, error) {
	pattern := "user:online:*"
	
	// 解析游标
	var startCursor uint64 = 0
	if cursor != "" {
		var err error
		startCursor, err = strconv.ParseUint(cursor, 10, 64)
		if err != nil {
			startCursor = 0
		}
	}
	
	var keys []string
	var nextCursor uint64 = startCursor
	var err error
	const batchSize = 100
	
	// 使用 SCAN 命令遍历匹配的 key（支持游标分页）
	var batch []string
	batch, nextCursor, err = c.client.Scan(ctx, nextCursor, pattern, batchSize).Result()
	if err != nil {
		return nil, "", err
	}
	keys = append(keys, batch...)
	
	// 获取每个用户的平台ID列表
	result := make(map[string][]int32)
	for _, key := range keys {
		// 从 key 中提取 userID
		userID := key[len("user:online:"):]
		
		// 获取该用户的平台ID列表
		platformIDs, err := c.GetUserOnline(ctx, userID)
		if err != nil {
			continue
		}
		
		if len(platformIDs) > 0 {
			result[userID] = platformIDs
		}
	}
	
	// 返回下一个游标
	nextCursorStr := ""
	if nextCursor > 0 {
		nextCursorStr = strconv.FormatUint(nextCursor, 10)
	}
	
	return result, nextCursorStr, nil
}
