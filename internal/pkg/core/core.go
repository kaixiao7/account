package core

import (
	"encoding/json"
	"net/http"

	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/validatetrans"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	switch typed := err.(type) {
	// json反序列化时类型错误
	case *json.UnmarshalTypeError:
		c.JSON(errno.ErrValidation.Http, RespResult{
			Code: errno.ErrValidation.Code,
			Msg:  errno.ErrValidation.Message,
		})
		return
	// gin验证错误
	case validator.ValidationErrors:
		// 参数验证失败错误，返回具体的验证信息
		c.JSON(errno.ErrValidation.Http, RespResult{
			Code: errno.ErrValidation.Code,
			Msg:  errno.ErrValidation.Message,
			Data: validatetrans.Translate(typed),
		})
		return
	case *errno.Err:
		// 自定义错误类型，返回自定义信息
		c.JSON(typed.Errno.Http, RespResult{
			Code: typed.Errno.Code,
			Msg:  typed.Errno.Message,
		})
		return
	default:
		// 其他错误，返回服务器异常
		c.JSON(http.StatusInternalServerError, RespResult{
			Code: errno.InternalServerError.Code,
			Msg:  errno.InternalServerError.Message,
		})
	}
}
