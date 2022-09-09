package middleware

import (
	"kaixiao7/account/internal/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestId 注入"X-Request-ID"到context和req/resp的header中
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get(constant.XRequestIDKey)

		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Set(constant.XRequestIDKey, requestId)
		c.Writer.Header().Set(constant.XRequestIDKey, requestId)
		c.Next()
	}
}
