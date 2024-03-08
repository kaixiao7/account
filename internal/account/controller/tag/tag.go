package tag

import "kaixiao7/account/internal/account/service"

type CategoryTagController struct {
	tagSrv service.CategoryTagSrv
}

func NewCategoryTagController() *CategoryTagController {
	return &CategoryTagController{
		tagSrv: service.NewCategoryTagSrv(),
	}
}
