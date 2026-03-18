package minio

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	imageInfoKeyPrefix    = "MINIO:IMAGE:"
	thumbnailKeyPrefix    = "MINIO:THUMBNAIL:"
	cacheExpiration       = time.Hour * 24 * 7 // 7 days
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*ImageInfo, error)) (*ImageInfo, error) {
	cacheKey := imageInfoKeyPrefix + key
	
	// Try to get from cache
	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var info ImageInfo
		if err := json.Unmarshal(data, &info); err == nil {
			return &info, nil
		}
	}
	
	// Get from source
	info, err := fn(ctx)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	if data, err := json.Marshal(info); err == nil {
		c.client.Set(ctx, cacheKey, data, cacheExpiration)
	}
	
	return info, nil
}

func (c *RedisCache) GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error) {
	cacheKey := fmt.Sprintf("%s%s:w%d:h%d:%s", thumbnailKeyPrefix, format, width, height, key)
	
	// Try to get from cache
	result, err := c.client.Get(ctx, cacheKey).Result()
	if err == nil {
		return result, nil
	}
	
	// Get from source
	thumbnailKey, err := minioCache(ctx)
	if err != nil {
		return "", err
	}
	
	// Cache the result
	c.client.Set(ctx, cacheKey, thumbnailKey, cacheExpiration)
	
	return thumbnailKey, nil
}

func (c *RedisCache) DelObjectImageInfoKey(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		cacheKeys[i] = imageInfoKeyPrefix + key
	}
	
	return c.client.Del(ctx, cacheKeys...).Err()
}

func (c *RedisCache) DelImageThumbnailKey(ctx context.Context, key string, format string, width int, height int) error {
	cacheKey := fmt.Sprintf("%s%s:w%d:h%d:%s", thumbnailKeyPrefix, format, width, height, key)
	return c.client.Del(ctx, cacheKey).Err()
}
