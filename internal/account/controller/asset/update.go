package asset

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AssetController) Update(c *gin.Context) {
	userId := controller.GetUserId(c)

	assetId, exist := controller.GetIntParamFromUrl(c, "assetId")
	if !exist {
		return
	}

	var asset model.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	asset.Id = assetId
	asset.UserId = userId
	asset.UpdateTime = time.Now().Unix()

	if err := a.assetSrv.Update(c, &asset); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
