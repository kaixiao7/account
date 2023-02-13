package budget

import (
	"github.com/gin-gonic/gin"
)

// Get 获取账本的总预算
// 其实这里有问题，这个表示的是获取账本下的所有预算，但是现在仅获取了总预算，并且应该用list表示
func (b *BudgetController) Get(c *gin.Context) {
	// bookId, err := strconv.Atoi(c.Param("bookId"))
	// if err != nil {
	// 	core.WriteRespErr(c, errno.New(errno.ErrPageNotFound))
	// 	return
	// }
	//
	// budget, err := b.budgetSrv.QueryBudget(c, bookId)
	// if err != nil {
	// 	core.WriteRespErr(c, err)
	// 	return
	// }
	//
	// core.WriteRespSuccess(c, budget)
}
