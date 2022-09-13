package tag

import "kaixiao7/account/internal/account/service"

type TagController struct {
	billSrv service.BillSrv
}

func NewTagController() *TagController {
	return &TagController{
		billSrv: service.NewBillSrv(),
	}
}
