package account

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AccountController) Push(c *gin.Context) {

	var account []*model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := a.accountSrv.Push(c, account, syncTime); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, controller.PushRes{
		SyncTime: syncTime,
	})
}
