package borrowlend

import "kaixiao7/account/internal/account/service"

type BorrowLendController struct {
	borrowLendSrv  service.BorrowLendSrv
	accountFlowSrv service.AccountFlowSrv
	accountSrv     service.AccountSrv
}

func NewBorrowController() *BorrowLendController {
	return &BorrowLendController{
		borrowLendSrv:  service.NewBorrowLendSrv(),
		accountFlowSrv: service.NewAccountFlowSrv(),
		accountSrv:     service.NewAccountSrv(),
	}
}
