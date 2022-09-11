package bill

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BillController) Delete(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId, ok := controller.GetIntParamFromUrl(c, "bookId")
	if !ok {
		return
	}

	billId, ok := controller.GetIntParamFromUrl(c, "billId")
	if !ok {
		return
	}

	err := b.billSrv.Delete(c, billId, userId, bookId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
