package middleware

import (
	"fmt"

	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/token"

	"github.com/gin-gonic/gin"
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

		uri := c.Request.RequestURI

		var identity int64
		var err error
		if uri == "/refresh" {
			identity, err = token.DecodeRefreshToken(t)
		} else {
			identity, err = token.DecodeAccessToken(t)
		}
		if err != nil {
			log.Error("token验证失败: " + err.Error())

			core.WriteRespErr(c, errno.New(errno.ErrTokenInvalid))
			c.Abort()
			return
		}

		c.Set(constant.XUserIdKey, identity)
		c.Next()
	}
}
