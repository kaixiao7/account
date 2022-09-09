package category

import "kaixiao7/account/internal/account/service"

type CategoryController struct {
	categorySrv service.CategorySrv
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		categorySrv: service.NewCategorySrv(),
	}
}
