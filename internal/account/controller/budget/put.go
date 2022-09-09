package budget

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

type budgetRequest struct {
	Budget float32 `json:"budget"`
}

// Put 更新账本预算
func (b *BudgetController) Put(c *gin.Context) {
	// 从body中获取json请求参数
	var budgetReq budgetRequest
	if err := c.ShouldBindJSON(&budgetReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	// 从请求路径中获取参数
	budgetId, err := strconv.Atoi(c.Param("budgetId"))
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrPageNotFound))
		return
	}

	if err := b.budgetSrv.SetBudget(c, budgetId, controller.GetUserId(c), budgetReq.Budget); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
