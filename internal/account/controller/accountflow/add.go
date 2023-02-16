package accountflow

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

type accountFlowReq struct {
	accountFlowUpdateReq
}

func (af *AccountFlowController) Add(c *gin.Context) {
	userId := controller.GetUserId(c)
	accountId, ok := controller.GetInt64ParamFromUrl(c, "accountId")
	if !ok {
		return
	}

	var accountFlowReq accountFlowReq
	if err := c.ShouldBindJSON(&accountFlowReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	accountFlow := model.AccountFlow{
		UserId:          userId,
		AccountId:       accountId,
		Type:            accountFlowReq.Type,
		Cost:            accountFlowReq.Cost,
		RecordTime:      accountFlowReq.RecordTime.Timestamp(),
		Remark:          accountFlowReq.Remark,
		TargetAccountId: &accountFlowReq.TargetAccountId,
		AssociateName:   accountFlowReq.AssociateName,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
	}

	if err := af.accountFlowSrv.Add(c, &accountFlow); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
