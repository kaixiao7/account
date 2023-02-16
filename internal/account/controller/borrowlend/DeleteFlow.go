package borrowlend

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// DeleteFlow 删除借入借出流水（还款、收款）
func (b *BorrowLendController) DeleteFlow(c *gin.Context) {
	userId := controller.GetUserId(c)

	flowId, ok := controller.GetInt64ParamFromUrl(c, "flowId")
	if !ok {
		return
	}

	if err := b.borrowLendSrv.DeleteBorrowLendFlow(c, flowId, userId); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
