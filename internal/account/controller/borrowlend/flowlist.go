package borrowlend

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// FlowList 查询借入借出流水列表
func (b *BorrowLendController) FlowList(c *gin.Context) {
	userId := controller.GetUserId(c)

	accountFlowId, ok := controller.GetInt64ParamFromUrl(c, "accountFlowId")
	if !ok {
		return
	}

	list, err := b.borrowLendSrv.QueryBorrowLendFlowList(c, accountFlowId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, list)
}
