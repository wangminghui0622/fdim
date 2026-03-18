package mw

import (
	"net/http"
	"strings"

	"fdim/pkg/errs"
	"fdim/pkg/tokenverify"
	"fdim/protocol/constant"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// CorsHandler gin cross-domain configuration.
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar",
		)
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("content-type", "application/json")
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GinParseOperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "header must have operationID"})
				c.Abort()
				return
			}
			c.Set(constant.OperationID, operationID)
		}
		c.Next()
	}
}

func GinParseToken(secretKey jwt.Keyfunc, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			for _, wApi := range whitelist {
				if strings.HasPrefix(c.Request.URL.Path, wApi) {
					c.Next()
					return
				}
			}

			token := c.Request.Header.Get(constant.Token)
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "header must have token"})
				c.Abort()
				return
			}

			claims, err := tokenverify.GetClaimFromToken(token, secretKey)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"errCode": 401, "errMsg": "invalid token"})
				c.Abort()
				return
			}

			c.Set(constant.OpUserPlatform, constant.PlatformIDToName(claims.PlatformID))
			c.Set(constant.OpUserID, claims.UserID)
			c.Next()
		}
	}
}

func CreateToken(userID string, accessSecret string, accessExpire int64, platformID int) (string, error) {
	claims := tokenverify.BuildClaims(userID, platformID, accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", errs.WrapMsg(err, "token.SignedString")
	}
	return tokenString, nil
}

func GinPanicErr(c *gin.Context, err any) {
	c.AbortWithStatus(http.StatusInternalServerError)
}
