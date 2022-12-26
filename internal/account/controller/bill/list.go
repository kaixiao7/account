package bill

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (b *BillController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId, ok := controller.GetIntParamFromUrl(c, "bookId")
	if !ok {
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 100
	}
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		pageNum = 1
	}

	ret, err := b.billSrv.QueryByPage(c, bookId, userId, pageSize, pageNum)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	core.WriteRespSuccess(c, ret)
}
