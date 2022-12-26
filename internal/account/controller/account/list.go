package account

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AccountController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	accounts, err := a.accountSrv.QueryByUserId(c, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, accounts)
}
