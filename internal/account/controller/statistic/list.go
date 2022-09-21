package statistic

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

func (s *StatisticController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}
	beginTime, err := timex.Parse(c.Query("beginTime"), timex.DatePattern)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}
	endTime, err := timex.Parse(c.Query("endTime"), timex.DatePattern)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	bills, err := s.statisticsSrv.QueryBill(c, bookId, userId, beginTime.Unix(), endTime.Unix())
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, bills)
}
