package asset

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (a *AssetController) Delete(c *gin.Context) {

	userId := controller.GetUserId(c)

	assetId, exist := controller.GetIntParamFromUrl(c, "assetId")
	if !exist {
		return
	}

	err := a.assetSrv.Delete(c, assetId, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
