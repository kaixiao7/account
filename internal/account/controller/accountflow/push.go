package accountflow

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (af *AccountFlowController) Push(c *gin.Context) {

	var flows []*model.AccountFlow
	if err := c.ShouldBindJSON(&flows); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := af.accountFlowSrv.Push(c, flows, syncTime); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, controller.PushRes{
		SyncTime: syncTime,
	})
}
