package errno

import "net/http"

var (
	OK = &Errno{Code: 0, Message: "SUCCESS", Http: http.StatusOK}

	InternalServerError = &Errno{Code: -1, Message: "服务器内部错误", Http: http.StatusInternalServerError}

	ErrValidation     = &Errno{Code: 10001, Message: "参数校验失败"}
	ErrPageNotFound   = &Errno{Code: 10002, Message: "请求地址不存在", Http: http.StatusNotFound}
	ErrToken          = &Errno{Code: 10003, Message: "JWT签名失败", Http: http.StatusInternalServerError}
	ErrIllegalOperate = &Errno{Code: 10004, Message: "非法操作"}

	ErrUserNotFound      = &Errno{Code: 10101, Message: "用户不存在"}
	ErrUserAlreadyExist  = &Errno{Code: 10102, Message: "用户已存在"}
	ErrTokenInvalid      = &Errno{Code: 10103, Message: "Token过期", Http: http.StatusUnauthorized}
	ErrPasswordIncorrect = &Errno{Code: 10104, Message: "用户名或密码错误"}

	ErrCategoryNotFound = &Errno{Code: 10201, Message: "分类不存在"}

	ErrBookNotFound = &Errno{Code: 10301, Message: "账本不存在"}

	ErrBillNotFound  = &Errno{Code: 10401, Message: "账单不存在"}
	ErrBillNotMore   = &Errno{Code: 10402, Message: "没有更多了"}
	ErrBillNotDelete = &Errno{Code: 10403, Message: "不允许删除他人账单"}
	ErrBillNotModify = &Errno{Code: 10403, Message: "不允许修改他人账单"}

	ErrAssetNotFound = &Errno{Code: 10501, Message: "账户不存在"}

	ErrAssetFlowNotFound     = &Errno{Code: 10601, Message: "流水不存在"}
	ErrAssetFlowAssociateNil = &Errno{Code: 10602, Message: "对方名称不能为空"}
)
