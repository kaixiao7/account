package asset

import "kaixiao7/account/internal/account/service"

type AssetController struct {
	assetSrv service.AssetSrv
}

func NewAssetController() *AssetController {
	return &AssetController{
		assetSrv: service.NewAssetSrv(),
	}
}
