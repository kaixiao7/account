package accountflow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (af *AccountFlowController) PullBook(c *gin.Context) {

	lastSyncTime, exist := controller.GetInt64ParamFromParam(c, "lastSyncTime")
	if !exist {
		return
	}

	bookId, exist := controller.GetInt64ParamFromUrl(c, "bookId")
	if !exist {
		return
	}

	startTime, exist := controller.GetInt64ParamFromParam(c, "startTime")
	if !exist {
		return
	}

	endTime, exist := controller.GetInt64ParamFromParam(c, "endTime")
	if !exist {
		return
	}

	// pageNum, exist := controller.GetInt64ParamFromParam(c, "pageNum")
	// if !exist {
	// 	return
	// }

	// flowPages, err := af.accountFlowSrv.PullBook(c, bookId, lastSyncTime, int(pageNum))

	flows, err := af.accountFlowSrv.PullBookRange(c, bookId, startTime, endTime, lastSyncTime)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, flows)
}
