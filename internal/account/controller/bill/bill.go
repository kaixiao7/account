package bill

import "kaixiao7/account/internal/account/service"

type BillController struct {
	billSrv service.BillSrv
}

func NewBillController() *BillController {
	return &BillController{billSrv: service.NewBillSrv()}
}
