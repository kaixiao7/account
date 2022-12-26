package accountflow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (af *AccountFlowController) List(c *gin.Context) {

	userId := controller.GetUserId(c)

	accountId, ok := controller.GetIntParamFromUrl(c, "accountId")
	if !ok {
		return
	}

	flows, err := af.accountFlowSrv.QueryByAccountId(c, accountId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, flows)
}
