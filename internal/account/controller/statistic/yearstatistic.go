package statistic

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// YearStatistic 年统计
func (s *StatisticController) YearStatistic(c *gin.Context) {
	userId := controller.GetUserId(c)

	year, ok := controller.GetIntParamFromUrl(c, "year")
	if !ok {
		return
	}
	bookId, err := strconv.Atoi(c.Query("bookId"))
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	ret, err := s.statisticsSrv.QueryByYear(c, bookId, userId, year)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, &ret)
}
