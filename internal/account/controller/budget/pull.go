package budget

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BudgetController) Pull(c *gin.Context) {
	bookId, exist := controller.GetInt64ParamFromUrl(c, "bookId")
	if !exist {
		return
	}

	lastSyncTime, exist := controller.GetInt64ParamFromParam(c, "lastSyncTime")
	if !exist {
		return
	}

	accounts, err := b.budgetSrv.Pull(c, bookId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, accounts)
}
