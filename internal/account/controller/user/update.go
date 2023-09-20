package user

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (u *UserController) Update(c *gin.Context) {

	userId := controller.GetUserId(c)

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	user.Id = userId
	user.UpdateTime = time.Now().Unix()

	if err := u.userSrv.Update(c, &user); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
