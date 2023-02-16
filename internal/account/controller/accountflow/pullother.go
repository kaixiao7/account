package accountflow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (af *AccountFlowController) PullOther(c *gin.Context) {

	userId := controller.GetUserId(c)

	lastSyncTime, exist := controller.GetInt64ParamFromParam(c, "lastSyncTime")
	if !exist {
		return
	}

	flows, err := af.accountFlowSrv.PullOther(c, userId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, flows)
}
