package statistic

import (
	"kaixiao7/account/internal/account/service"
)

type StatisticController struct {
	statisticsSrv service.StatisticsSrv
}

func NewStatisticController() *StatisticController {
	return &StatisticController{
		statisticsSrv: service.NewStatisticsSrv(),
	}
}
