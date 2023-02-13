package category

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (c *CategoryController) Pull(ctx *gin.Context) {

	bookId, exist := controller.GetIntParamFromUrl(ctx, "bookId")
	if !exist {
		return
	}

	lastSyncTime, exist := controller.GetInt64ParamFromParam(ctx, "lastSyncTime")
	if !exist {
		return
	}

	categories, err := c.categorySrv.Pull(ctx, bookId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	core.WriteRespSuccess(ctx, categories)
}
