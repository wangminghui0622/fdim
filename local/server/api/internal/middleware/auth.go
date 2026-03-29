package middleware

import (
	"fmt"
	"net/http"
	"strings"
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

var whitelist = []string{
	"/auth/get_admin_token",
	"/auth/parse_token",
	"/auth/get_user_token",
	"/account/login",
	"/account/register",
	"/account/change_password",
	"/account/send_verify_code",
	"/account/verify_code",
	"/account/reset_password",
}

func isWhitelisted(path string) bool {
	for _, item := range whitelist {
		if strings.HasPrefix(path, item) {
			return true
		}
	}
	return false
}

// AuthMiddleware 认证中间件，?token 中提取用户信息并设置?context
func AuthMiddleware(secret string) rest.Middleware {
	fmt.Printf("[AuthMiddleware] Initialized with secret length: %d\n", len(secret))
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			operationID := r.Header.Get("operationID")
			if operationID == "" {
				operationID = r.Header.Get("OperationID")
			}
			if operationID != "" {
				ctx = mcontext.SetOperationID(ctx, operationID)
			}
			if mcontext.GetOperationID(ctx) == "" {
				ctx = mcontext.SetOperationID(ctx, time.Now().Format("20060102150405"))
			}

			if isWhitelisted(r.URL.Path) {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			tokenStr := r.Header.Get("token")
			if tokenStr == "" {
				tokenStr = r.Header.Get("Token")
			}
			fmt.Printf("[AuthMiddleware] Request: %s, Token present: %v, Token length: %d\n", r.URL.Path, tokenStr != "", len(tokenStr))

			if tokenStr == "" {
				http.Error(w, "header must have token", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				fmt.Printf("[AuthMiddleware] Token parse error: %v\n", err)
				http.Error(w, "token parse failed", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				fmt.Printf("[AuthMiddleware] Token invalid\n")
				http.Error(w, "token invalid", http.StatusUnauthorized)
				return
			}

			c, ok := token.Claims.(*claims)
			if !ok || c.UserID == "" {
				fmt.Printf("[AuthMiddleware] Token claims invalid or userID empty\n")
				http.Error(w, "token invalid", http.StatusUnauthorized)
				return
			}

			fmt.Printf("[AuthMiddleware] Token parsed, UserID: %s\n", c.UserID)
			ctx = mcontext.SetOpUserID(ctx, c.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
