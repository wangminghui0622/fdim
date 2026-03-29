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

// SetUserOnline û״̬ƽ̨
func (c *UserCache) SetUserOnline(ctx context.Context, userID string, platformIDs []int32) error {
	_ = "user:online:" + userID
	// ʵ߼ platformIDs 洢 Redis Set
	// ʵ֣ʵӦʹ Redis Set 洢ƽ̨ID
	return nil
}

// SetUserOnlineStatus û״̬֧ߺƽ̨б
func (c *UserCache) SetUserOnlineStatus(ctx context.Context, userID string, online []int32, offline []int32) error {
	key := "user:online:" + userID
	
	// ƽ̨
	for _, platformID := range online {
		member := strconv.FormatInt(int64(platformID), 10)
		if err := c.client.SAdd(ctx, key, member).Err(); err != nil {
			return err
		}
	}
	
	// Ƴƽ̨
	for _, platformID := range offline {
		member := strconv.FormatInt(int64(platformID), 10)
		if err := c.client.SRem(ctx, key, member).Err(); err != nil {
			return err
		}
	}
	
	// ùʱ䣨4Сʱ
	if err := c.client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return err
	}
	
	return nil
}

// GetUserOnline ȡû״̬
func (c *UserCache) GetUserOnline(ctx context.Context, userID string) ([]int32, error) {
	key := "user:online:" + userID
	
	//  Redis Set ȡƽ̨ID
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

// GetAllOnlineUsers ȡû֧αҳ
func (c *UserCache) GetAllOnlineUsers(ctx context.Context, cursor string) (map[string][]int32, string, error) {
	pattern := "user:online:*"
	
	// α
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
	
	// ʹ SCAN ƥ key֧αҳ
	var batch []string
	batch, nextCursor, err = c.client.Scan(ctx, nextCursor, pattern, batchSize).Result()
	if err != nil {
		return nil, "", err
	}
	keys = append(keys, batch...)
	
	// ȡÿûƽ̨IDб
	result := make(map[string][]int32)
	for _, key := range keys {
		//  key ȡ userID
		userID := key[len("user:online:"):]
		
		// ȡûƽ̨IDб
		platformIDs, err := c.GetUserOnline(ctx, userID)
		if err != nil {
			continue
		}
		
		if len(platformIDs) > 0 {
			result[userID] = platformIDs
		}
	}
	
	// һα
	nextCursorStr := ""
	if nextCursor > 0 {
		nextCursorStr = strconv.FormatUint(nextCursor, 10)
	}
	
	return result, nextCursorStr, nil
}
