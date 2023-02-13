package budget

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *BudgetController) Push(c *gin.Context) {

	var budgets []*model.Budget
	if err := c.ShouldBindJSON(&budgets); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := a.budgetSrv.Push(c, budgets, syncTime); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, controller.PushRes{
		SyncTime: syncTime,
	})
}
