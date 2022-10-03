package assetflow

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

type assetFlowUpdateReq struct {
	Type          int            `db:"type" json:"type"  binding:"required,numeric"`
	Cost          float64        `db:"cost" json:"cost"  binding:"required,numeric"`
	RecordTime    timex.JsonTime `db:"record_time" json:"record_time"  binding:"required"`
	Remark        string         `db:"remark" json:"remark,omitempty"`
	CategoryId    int            `db:"category_id" json:"category_id,omitempty"`
	TargetAssetId int            `db:"target_asset_id" json:"target_asset_id,omitempty"`
	AssociateName string         `db:"associate_name" json:"associate_name,omitempty"`
}

func (af *AssetFlowController) Update(c *gin.Context) {
	userId := controller.GetUserId(c)
	assetId, ok := controller.GetIntParamFromUrl(c, "assetId")
	if !ok {
		return
	}
	assetFlowId, ok := controller.GetIntParamFromUrl(c, "assetFlowId")
	if !ok {
		return
	}

	var assetFlowReq assetFlowUpdateReq
	if err := c.ShouldBindJSON(&assetFlowReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	assetFlow := model.AssetFlow{
		Id:            assetFlowId,
		UserId:        userId,
		AssetId:       assetId,
		Type:          assetFlowReq.Type,
		Cost:          assetFlowReq.Cost,
		RecordTime:    assetFlowReq.RecordTime.Timestamp(),
		Remark:        assetFlowReq.Remark,
		CategoryId:    &assetFlowReq.CategoryId,
		TargetAssetId: &assetFlowReq.TargetAssetId,
		AssociateName: assetFlowReq.AssociateName,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	if err := af.assetFlowSrv.Update(c, &assetFlow); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
