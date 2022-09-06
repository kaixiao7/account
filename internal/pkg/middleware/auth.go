package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/token"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			log.Error("token长度为0")
			core.WriteRespErr(c, errno.New(errno.ErrTokenInvalid))
			c.Abort()
			return
		}

		var t string
		// Parse the header to get the token part.
		fmt.Sscanf(header, "Bearer %s", &t)

		userId, err := token.Parse(t)
		if err != nil {
			log.Error("token验证失败")
			core.WriteRespErr(c, errno.New(errno.ErrTokenInvalid))
			c.Abort()
			return
		}

		c.Set(constant.XUserIdKey, userId)
		c.Next()
	}
}
