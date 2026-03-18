package ctxutil

import (
	"context"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	tokenKey     contextKey = "token"
	platformKey  contextKey = "platform"
	operationKey contextKey = "operationID"
)

// GetUserID 从 context 中获取用户ID
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(userIDKey); v != nil {
		if userID, ok := v.(string); ok {
			return userID
		}
	}
	return ""
}

// SetUserID 设置用户ID到 context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetToken 从 context 中获取 token
func GetToken(ctx context.Context) string {
	if v := ctx.Value(tokenKey); v != nil {
		if token, ok := v.(string); ok {
			return token
		}
	}
	return ""
}

// SetToken 设置 token 到 context
func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// GetPlatform 从 context 中获取平台
func GetPlatform(ctx context.Context) int {
	if v := ctx.Value(platformKey); v != nil {
		if platform, ok := v.(int); ok {
			return platform
		}
	}
	return 0
}

// SetPlatform 设置平台到 context
func SetPlatform(ctx context.Context, platform int) context.Context {
	return context.WithValue(ctx, platformKey, platform)
}

// GetOperationID 从 context 中获取操作ID
func GetOperationID(ctx context.Context) string {
	if v := ctx.Value(operationKey); v != nil {
		if opID, ok := v.(string); ok {
			return opID
		}
	}
	return ""
}

// SetOperationID 设置操作ID到 context
func SetOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, operationKey, operationID)
}
