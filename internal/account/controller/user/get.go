package user

import (
	"github.com/gin-gonic/gin"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
)

func (u *UserController) Get(c *gin.Context) {

	userId, exist := c.Get(constant.XUserIdKey)

	if !exist {
		core.WriteRespErr(c, errno.New(errno.ErrToken))
		return
	}

	user, err := u.userSrv.GetById(userId.(int))
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, user)
}
