package borrowlend

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (b *BorrowLendController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	borrowType := c.Query("type")
	borrowTypeNum, err := strconv.Atoi(borrowType)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	list, err := b.borrowLendSrv.QueryBorrowLendList(c, userId, borrowTypeNum)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, list)
}
