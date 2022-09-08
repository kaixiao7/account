package middleware

import (
	"kaixiao7/account/internal/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RequestId 注入"X-Request-ID"到context和req/resp的header中
func SqlDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.SqlDBKey, db)
		c.Next()
	}
}
