package controller

import (
	"strconv"

	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// GetUserId 从上下文中获取认证用户id
func GetUserId(c *gin.Context) int {
	userId, exist := c.Get(constant.XUserIdKey)

	if !exist {
		core.WriteRespErr(c, errno.New(errno.InternalServerError))
		return 0
	}

	return userId.(int)
}

// GetIntParamFromUrl 从请求路径中获取指定的int类型参数
func GetIntParamFromUrl(c *gin.Context, paramName string) (int, bool) {
	bookId, err := strconv.Atoi(c.Param(paramName))
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrPageNotFound))
		return 0, false
	}

	return bookId, true
}
