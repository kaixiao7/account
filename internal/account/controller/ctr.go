package controller

import (
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) int {
	userId, exist := c.Get(constant.XUserIdKey)

	if !exist {
		core.WriteRespErr(c, errno.New(errno.InternalServerError))
		return 0
	}

	return userId.(int)
}
