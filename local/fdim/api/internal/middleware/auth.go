package middleware

import (
	"fmt"
	"net/http"
	"time"

	"fdim/pkg/mcontext"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest"
)

// claims 与 tokenverify 中的 claims 结构保持一致
type claims struct {
	UserID     string
	UserType   int32
	PlatformID int32
	jwt.RegisteredClaims
}

// AuthMiddleware 认证中间件，从 token 中提取用户信息并设置到 context
func AuthMiddleware(secret string) rest.Middleware {
	fmt.Printf("[AuthMiddleware] Initialized with secret length: %d\n", len(secret))
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			
			// 从 header 中获取 token
			tokenStr := r.Header.Get("token")
			if tokenStr == "" {
				tokenStr = r.Header.Get("Token")
			}
			fmt.Printf("[AuthMiddleware] Request: %s, Token present: %v, Token length: %d\n", r.URL.Path, tokenStr != "", len(tokenStr))
			
			// 从 header 中获取 operationID
			operationID := r.Header.Get("operationID")
			if operationID == "" {
				operationID = r.Header.Get("OperationID")
			}
			if operationID != "" {
				ctx = mcontext.SetOperationID(ctx, operationID)
			}
			
			// 设置默认 operationID
			if mcontext.GetOperationID(ctx) == "" {
				ctx = mcontext.SetOperationID(ctx, time.Now().Format("20060102150405"))
			}
			
			// 解析 token 获取用户信息
			if tokenStr != "" {
				token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				})
				if err != nil {
					fmt.Printf("[AuthMiddleware] Token parse error: %v\n", err)
				} else if token.Valid {
					if c, ok := token.Claims.(*claims); ok {
						fmt.Printf("[AuthMiddleware] Token parsed, UserID: %s\n", c.UserID)
						ctx = mcontext.SetOpUserID(ctx, c.UserID)
					}
				} else {
					fmt.Printf("[AuthMiddleware] Token invalid\n")
				}
			}
			
			// 使用更新后的 context 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
