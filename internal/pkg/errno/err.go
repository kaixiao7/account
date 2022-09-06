package errno

import (
	"fmt"
	"net/http"
)

// Errno 错误码类型
type Errno struct {
	Code    int
	Http    int
	Message string
}

// Err 自定义错误类型（带有code的错误）
type Err struct {
	Errno *Errno
	Err   error
}

func NewWithError(errno *Errno, err error) *Err {
	if errno.Http == 0 {
		// http状态码默认值为400
		errno.Http = http.StatusBadRequest
	}
	return &Err{Errno: errno, Err: err}
}

func New(errno *Errno) *Err {
	return &Err{Errno: errno}
}

func (e Err) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Err - code: %d", e.Errno.Code)
	}

	return fmt.Sprintf("Err - code: %d, error: %s", e.Errno.Code, e.Err)
}
