package errno

import "net/http"

var (
	OK = &Errno{Code: 0, Message: "SUCCESS", Http: http.StatusOK}

	InternalServerError = &Errno{Code: -1, Message: "服务器内部错误", Http: http.StatusInternalServerError}

	ErrValidation   = &Errno{Code: 10001, Message: "参数校验失败"}
	ErrPageNotFound = &Errno{Code: 10002, Message: "请求地址不存在", Http: http.StatusNotFound}
	ErrToken        = &Errno{Code: 10003, Message: "JWT签名失败", Http: http.StatusInternalServerError}

	ErrUserNotFound      = &Errno{Code: 10101, Message: "用户不存在"}
	ErrUserAlreadyExist  = &Errno{Code: 10102, Message: "用户已存在"}
	ErrTokenInvalid      = &Errno{Code: 10103, Message: "Token过期", Http: http.StatusUnauthorized}
	ErrPasswordIncorrect = &Errno{Code: 10104, Message: "用户名或密码错误"}
)
