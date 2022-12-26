package accountflow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (af *AccountFlowController) Delete(c *gin.Context) {
	userId := controller.GetUserId(c)
	accountFlowId, ok := controller.GetIntParamFromUrl(c, "accountFlowId")
	if !ok {
		return
	}

	err := af.accountFlowSrv.Delete(c, accountFlowId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
