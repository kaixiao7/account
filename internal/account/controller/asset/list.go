package asset

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AssetController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	assets, err := a.assetSrv.QueryByUserId(c, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, assets)
}
