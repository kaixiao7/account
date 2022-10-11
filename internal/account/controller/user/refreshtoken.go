package user

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (u *UserController) RefreshToken(c *gin.Context) {
	userId := controller.GetUserId(c)

	resp, err := generateTokens(userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, resp)
}
