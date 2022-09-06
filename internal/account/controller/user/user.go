package user

import "kaixiao7/account/internal/account/service"

type UserController struct {
	userSrv service.UserSrv
}

// NewUserController 创建用户处理器
func NewUserController() *UserController {
	return &UserController{userSrv: service.NewUserSrv()}
}
