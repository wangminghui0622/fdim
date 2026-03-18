package rpccache

import (
	"context"
	"time"

	"fdim/pkg/config"
	"fdim/pkg/localcache"
	"fdim/pkg/log"
	"fdim/pkg/rpcli"
	"fdim/pkg/storage/cache/cachekey"
	"fdim/protocol/auth"
	"github.com/redis/go-redis/v9"
)

func NewAuthLocalCache(client *rpcli.AuthClient, localCache *config.LocalCache, cli redis.UniversalClient) *AuthLocalCache {
	lc := localCache.Auth
	log.ZDebug(context.Background(), "AuthLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &AuthLocalCache{
		client: client,
		local: localcache.New[[]byte](
			localcache.WithLocalSlotNum(lc.SlotNum),
			localcache.WithLocalSlotSize(lc.SlotSize),
			localcache.WithLinkSlotNum(lc.SlotNum),
			localcache.WithLocalSuccessTTL(lc.Success()),
			localcache.WithLocalFailedTTL(lc.Failed()),
		),
	}
	if lc.Enable() {
		go subscriberRedisDeleteCache(context.Background(), cli, lc.Topic, x.local.DelLocal)
	}
	return x
}

type AuthLocalCache struct {
	client *rpcli.AuthClient
	local  localcache.Cache[[]byte]
}

func (a *AuthLocalCache) GetExistingToken(ctx context.Context, userID string, platformID int) (val map[string]int, err error) {
	resp, err := a.getExistingToken(ctx, userID, platformID)
	if err != nil {
		return nil, err
	}

	// Convert map[string]int32 to map[string]int
	res := make(map[string]int, len(resp.TokenStates))
	for k, v := range resp.TokenStates {
		res[k] = int(v)
	}

	return res, nil
}

func (a *AuthLocalCache) getExistingToken(ctx context.Context, userID string, platformID int) (val *auth.GetExistingTokenResp, err error) {
	start := time.Now()
	log.ZDebug(ctx, "AuthLocalCache GetExistingToken req", "userID", userID, "platformID", platformID)
	defer func() {
		if err != nil {
			log.ZError(ctx, "AuthLocalCache GetExistingToken error", err, "cost", time.Since(start), "userID", userID, "platformID", platformID)
		} else {
			log.ZDebug(ctx, "AuthLocalCache GetExistingToken resp", "cost", time.Since(start), "userID", userID, "platformID", platformID, "val", val)
		}
	}()

	var cache cacheProto[auth.GetExistingTokenResp]

	return cache.Unmarshal(a.local.Get(ctx, cachekey.GetTokenKey(userID, platformID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "AuthLocalCache GetExistingToken call rpc", "userID", userID, "platformID", platformID)
		return cache.Marshal(a.client.AuthClient.GetExistingToken(ctx, &auth.GetExistingTokenReq{UserID: userID, PlatformID: int32(platformID)}))
	}))
}
