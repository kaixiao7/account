package middleware

import (
	"kaixiao7/account/internal/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SqlDB 将sqlx.DB放入上下文中，供store使用
func SqlDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.SqlDBKey, db)
		c.Next()
	}
}
