package account

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AccountController) Update(c *gin.Context) {
	userId := controller.GetUserId(c)

	accountId, exist := controller.GetIntParamFromUrl(c, "accountId")
	if !exist {
		return
	}

	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	account.Id = accountId
	account.UserId = userId
	account.UpdateTime = time.Now().Unix()

	if err := a.accountSrv.Update(c, &account); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
