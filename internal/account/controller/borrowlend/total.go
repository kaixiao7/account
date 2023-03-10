package borrowlend

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BorrowLendController) Total(c *gin.Context) {
	userId := controller.GetUserId(c)

	total, err := b.borrowLendSrv.QueryTotal(c, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, total)
}
