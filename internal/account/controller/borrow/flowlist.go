package borrow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// FlowList 查询借入借出流水列表
func (b *BorrowController) FlowList(c *gin.Context) {
	userId := controller.GetUserId(c)

	assetFlowId, ok := controller.GetIntParamFromUrl(c, "assetFlowId")
	if !ok {
		return
	}

	list, err := b.borrowSrv.QueryBorrowFlowList(c, assetFlowId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, list)
}
