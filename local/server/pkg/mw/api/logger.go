package api

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinLogger 简化的日志中间件
func GinLogger(skipPath ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

type httpResponse struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (r *httpResponse) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

// GinError 返回错误响应
func GinError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
}
