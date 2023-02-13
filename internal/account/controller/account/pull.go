package account

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AccountController) Pull(c *gin.Context) {
	userId := controller.GetUserId(c)

	lastSyncTime, exist := controller.GetInt64ParamFromParam(c, "lastSyncTime")
	if !exist {
		return
	}

	accounts, err := a.accountSrv.Pull(c, userId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, accounts)
}
