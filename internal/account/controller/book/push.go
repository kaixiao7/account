package book

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BookController) Push(c *gin.Context) {

	var books []*model.Book
	if err := c.ShouldBindJSON(&books); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := b.bookSrv.Push(c, books, syncTime); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, controller.PushRes{
		SyncTime: syncTime,
	})
}
