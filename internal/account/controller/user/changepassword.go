package user

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

type PwdModel struct {
	OldPwd string `json:"old_pwd" binding:"required,max=100"`
	NewPwd string `json:"new_pwd" binding:"required,max=100"`
}

func (u *UserController) ChangePassword(c *gin.Context) {

	userId := controller.GetUserId(c)

	var pwd PwdModel

	if err := c.ShouldBindJSON(&pwd); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	if err := u.userSrv.ChangePassword(c, pwd.OldPwd, pwd.NewPwd, userId); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
