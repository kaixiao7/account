package category

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (c *CategoryController) Push(ctx *gin.Context) {

	var categories []*model.Category
	if err := ctx.ShouldBindJSON(&categories); err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := c.categorySrv.Push(ctx, categories, syncTime); err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	core.WriteRespSuccess(ctx, controller.PushRes{
		SyncTime: syncTime,
	})
}
