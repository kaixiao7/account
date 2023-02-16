package category

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (c *CategoryController) List(ctx *gin.Context) {
	bookId, ok := controller.GetInt64ParamFromUrl(ctx, "bookId")
	if !ok {
		return
	}

	categories, err := c.categorySrv.QueryAll(ctx, bookId)
	if err != nil {
		core.WriteRespErr(ctx, err)
		return
	}

	core.WriteRespSuccess(ctx, categories)
}
