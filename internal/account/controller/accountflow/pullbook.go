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

	bookId, exist := controller.GetIntParamFromUrl(c, "bookId")
	if !exist {
		return
	}

	pageNum, exist := controller.GetInt64ParamFromParam(c, "pageNum")
	if !exist {
		return
	}

	flowPages, err := af.accountFlowSrv.PullBook(c, int(bookId), lastSyncTime, int(pageNum))
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, flowPages)
}
