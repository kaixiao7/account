package user

import (
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (u *UserController) Get(c *gin.Context) {

	userId, exist := c.Get(constant.XUserIdKey)

	if !exist {
		core.WriteRespErr(c, errno.New(errno.ErrToken))
		return
	}

	user, err := u.userSrv.GetById(c, userId.(int))
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, user)
}
