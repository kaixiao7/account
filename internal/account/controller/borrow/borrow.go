package borrow

import "kaixiao7/account/internal/account/service"

type BorrowController struct {
	borrowSrv    service.BorrowSrv
	assetFlowSrv service.AssetFlowSrv
	assetSrv     service.AssetSrv
}

func NewBorrowController() *BorrowController {
	return &BorrowController{
		borrowSrv:    service.NewBorrowSrv(),
		assetFlowSrv: service.NewAssertFlowSrv(),
		assetSrv:     service.NewAssetSrv(),
	}
}
