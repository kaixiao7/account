package assetflow

import (
	"kaixiao7/account/internal/account/service"
)

type AssetFlowController struct {
	assetFlowSrv service.AssetFlowSrv
}

func NewAssetFlowController() *AssetFlowController {
	return &AssetFlowController{
		assetFlowSrv: service.NewAssertFlowSrv(),
	}
}
