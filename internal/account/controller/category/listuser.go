package category

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (c *CategoryController) ListUser(ctx *gin.Context) {
	userId := controller.GetUserId(ctx)

	categories, err := c.categorySrv.QueryByUserId(ctx, userId)
	if err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	core.WriteRespSuccess(ctx, categories)
}
