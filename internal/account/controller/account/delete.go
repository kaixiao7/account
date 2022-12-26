package account

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AccountController) Delete(c *gin.Context) {

	userId := controller.GetUserId(c)

	accountId, exist := controller.GetIntParamFromUrl(c, "accountId")
	if !exist {
		return
	}

	err := a.accountSrv.Delete(c, accountId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
