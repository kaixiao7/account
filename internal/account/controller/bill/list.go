package bill

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

func (b *BillController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId, ok := controller.GetIntParamFromUrl(c, "bookId")
	if !ok {
		return
	}

	date := c.Query("date")
	// 如果没有传递参数，则默认为当前时间
	if date == "" {
		date = timex.Format(time.Now(), timex.DatePattern)
	}
	// 时间转换
	parseDate, err := timex.Parse(date, timex.DatePattern)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	bills, err := b.billSrv.QueryByTime(c, bookId, userId, parseDate)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, bills)
}
