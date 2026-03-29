package database

import (
	"context"
	"fmt"

	"fdim/pkg/tokenverify"
	constantpb "fdim/protocol/constant"
	"github.com/redis/go-redis/v9"
)

type AuthDatabase interface {
	CreateToken(ctx context.Context, userID string, platformID int32, userType int32) (token string, expireSeconds int64, err error)
	GetTokens(ctx context.Context, userID string, platformID int32) (map[string]int32, error)
	SetTokenStatus(ctx context.Context, userID string, platformID int32, token string, status int32) error
	BatchKickTokens(ctx context.Context, userID string, platformID int32, tokens []string) error
}

type authDatabase struct {
	redis      *redis.Client
	token      *tokenverify.Token
	multiLogin bool
}

func NewAuthDatabase(rdb *redis.Client, token *tokenverify.Token, multiLogin bool) AuthDatabase {
	return &authDatabase{
		redis:      rdb,
		token:      token,
		multiLogin: multiLogin,
	}
}

func (a *authDatabase) tokenKey(userID string, platformID int32) string {
	return fmt.Sprintf("UID_PID_TOKEN_STATUS:%s:%d", userID, platformID)
}

func (a *authDatabase) CreateToken(ctx context.Context, userID string, platformID int32, userType int32) (string, int64, error) {
	if a.redis == nil || a.token == nil {
		return "", 0, fmt.Errorf("auth database not initialized")
	}

	tokenStr, expire, err := a.token.CreateToken(userID, userType)
	if err != nil {
		return "", 0, err
	}

	key := a.tokenKey(userID, platformID)
	existing, err := a.redis.HGetAll(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", 0, err
	}

	if !a.multiLogin && len(existing) > 0 {
		for tk := range existing {
			if tk == tokenStr {
				continue
			}
			a.redis.HSet(ctx, key, tk, constantpb.KickedToken)
		}
	}

	if err = a.redis.HSet(ctx, key, tokenStr, constantpb.NormalToken).Err(); err != nil {
		return "", 0, err
	}
	if err = a.redis.Expire(ctx, key, expire).Err(); err != nil {
		return "", 0, err
	}

	return tokenStr, int64(expire.Seconds()), nil
}

func (a *authDatabase) GetTokens(ctx context.Context, userID string, platformID int32) (map[string]int32, error) {
	if a.redis == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}
	key := a.tokenKey(userID, platformID)
	m, err := a.redis.HGetAll(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	result := make(map[string]int32, len(m))
	for tk, v := range m {
		switch v {
		case "0":
			result[tk] = constantpb.NormalToken
		case "2":
			result[tk] = constantpb.KickedToken
		case "3":
			result[tk] = constantpb.ExpiredToken
		default:
			result[tk] = constantpb.InValidToken
		}
	}
	return result, nil
}

func (a *authDatabase) SetTokenStatus(ctx context.Context, userID string, platformID int32, token string, status int32) error {
	if a.redis == nil {
		return fmt.Errorf("redis client not initialized")
	}
	key := a.tokenKey(userID, platformID)
	return a.redis.HSet(ctx, key, token, status).Err()
}

func (a *authDatabase) BatchKickTokens(ctx context.Context, userID string, platformID int32, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	if a.redis == nil {
		return fmt.Errorf("redis client not initialized")
	}
	key := a.tokenKey(userID, platformID)
	pipe := a.redis.TxPipeline()
	for _, tk := range tokens {
		pipe.HSet(ctx, key, tk, constantpb.KickedToken)
	}
	_, err := pipe.Exec(ctx)
	return err
}
