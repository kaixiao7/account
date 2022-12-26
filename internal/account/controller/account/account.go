package account

import "kaixiao7/account/internal/account/service"

type AccountController struct {
	accountSrv service.AccountSrv
}

func NewAccountController() *AccountController {
	return &AccountController{
		accountSrv: service.NewAccountSrv(),
	}
}
