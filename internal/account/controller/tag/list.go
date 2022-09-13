package tag

import (
	"strconv"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (t *TagController) List(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId := c.Query("bookId")
	if bookId == "" {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}
	bid, err := strconv.Atoi(bookId)
	if err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrValidation))
		return
	}

	tags, err := t.billSrv.QueryTag(c, bid, userId)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, tags)
}
