package assetflow

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

type assetFlowReq struct {
	AssetId int `db:"asset_id" json:"asset_id"  binding:"required,numeric"`
	assetFlowUpdateReq
}

func (af *AssetFlowController) Add(c *gin.Context) {
	userId := controller.GetUserId(c)

	var assetFlowReq assetFlowReq
	if err := c.ShouldBindJSON(&assetFlowReq); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	assetFlow := model.AssetFlow{
		UserId:        userId,
		AssetId:       assetFlowReq.AssetId,
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

	if err := af.assetFlowSrv.Add(c, &assetFlow); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}