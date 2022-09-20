package asset

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AssetController) Add(c *gin.Context) {
	userId := controller.GetUserId(c)

	var asset model.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	asset.UserId = userId
	asset.CreateTime = time.Now().Unix()
	asset.UpdateTime = time.Now().Unix()

	if err := a.assetSrv.Add(c, &asset); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
