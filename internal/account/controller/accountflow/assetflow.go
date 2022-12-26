package accountflow

import (
	"kaixiao7/account/internal/account/service"
)

type AccountFlowController struct {
	accountFlowSrv service.AccountFlowSrv
}

func NewAccountFlowController() *AccountFlowController {
	return &AccountFlowController{
		accountFlowSrv: service.NewAccountFlowSrv(),
	}
}
