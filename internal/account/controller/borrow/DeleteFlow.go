package borrow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BorrowController) DeleteFlow(c *gin.Context) {
	userId := controller.GetUserId(c)

	flowId, ok := controller.GetIntParamFromUrl(c, "flowId")
	if !ok {
		return
	}

	if err := b.borrowSrv.DeleteBorrowFlow(c, flowId, userId); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
