package core

import (
	"github.com/gin-gonic/gin"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"net/http"
)

type RespResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func WriteRespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, RespResult{
		Code: errno.OK.Code,
		Msg:  errno.OK.Message,
		Data: data,
	})
}

func WriteRespErr(c *gin.Context, err error) {
	log.Errorf("%+v", err)
	if en, ok := err.(*errno.Err); ok {
		c.JSON(en.Errno.Http, RespResult{
			Code: en.Errno.Code,
			Msg:  en.Errno.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, RespResult{
			Code: errno.InternalServerError.Code,
			Msg:  errno.InternalServerError.Message,
		})
	}
}
