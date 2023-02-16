package accountflow

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

type accountFlowUpdateReq struct {
	Type            int            `db:"type" json:"type"  binding:"required,numeric"`
	Cost            float64        `db:"cost" json:"cost"  binding:"required,numeric"`
	RecordTime      timex.JsonTime `db:"record_time" json:"record_time"  binding:"required"`
	Remark          string         `db:"remark" json:"remark,omitempty"`
	CategoryId      int64          `db:"category_id" json:"category_id,omitempty"`
	TargetAccountId int64          `db:"target_account_id" json:"target_account_id,omitempty"`
	AssociateName   string         `db:"associate_name" json:"associate_name,omitempty"`
}

func (af *AccountFlowController) Update(c *gin.Context) {
	userId := controller.GetUserId(c)
	accountId, ok := controller.GetInt64ParamFromUrl(c, "accountId")
	if !ok {
		return
	}
	accountFlowId, ok := controller.GetInt64ParamFromUrl(c, "accountFlowId")
	if !ok {
		return
	}

	var accountFlowReq accountFlowUpdateReq
	if err := c.ShouldBindJSON(&accountFlowReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	accountFlow := model.AccountFlow{
		Id:              accountFlowId,
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

	if err := af.accountFlowSrv.Update(c, &accountFlow); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
