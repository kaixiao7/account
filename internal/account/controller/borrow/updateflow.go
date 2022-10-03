package borrow

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

type updateFlowReq struct {
	AssetId    int            `json:"asset_id" binding:"required,numeric"`
	Cost       float64        `json:"cost" binding:"required,numeric"`
	RecordTime timex.JsonTime `json:"record_time" binding:"required"`
	Type       int            `json:"type,omitempty" binding:"required,numeric"`
	Remark     string         `json:"remark,omitempty"`
}

func (b *BorrowController) UpdateFlow(c *gin.Context) {
	userId := controller.GetUserId(c)

	assetFlowId, ok := controller.GetIntParamFromUrl(c, "assetFlowId")
	if !ok {
		return
	}
	flowId, ok := controller.GetIntParamFromUrl(c, "flowId")
	if !ok {
		return
	}

	var flowReq updateFlowReq
	if err := c.ShouldBindJSON(&flowReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	flow := model.BorrowFlow{
		Id:          flowId,
		AssetFlowId: assetFlowId,
		AssetId:     flowReq.AssetId,
		Cost:        flowReq.Cost,
		RecordTime:  flowReq.RecordTime.Timestamp(),
		Type:        flowReq.Type,
		Remark:      flowReq.Remark,
	}
	if err := b.borrowSrv.UpdateBorrowFlow(c, &flow, userId); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
