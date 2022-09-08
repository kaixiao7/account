package book

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// List 查询账本列表
func (b *BookController) List(c *gin.Context) {

	list, err := b.bookSrv.QueryBookList(c, controller.GetUserId(c))
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, list)
}
