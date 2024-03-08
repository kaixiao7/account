package tag

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (c *CategoryTagController) Pull(ctx *gin.Context) {

	bookId, exist := controller.GetInt64ParamFromUrl(ctx, "bookId")
	if !exist {
		return
	}

	lastSyncTime, exist := controller.GetInt64ParamFromParam(ctx, "lastSyncTime")
	if !exist {
		return
	}

	categories, err := c.tagSrv.Pull(ctx, bookId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	core.WriteRespSuccess(ctx, categories)
}
