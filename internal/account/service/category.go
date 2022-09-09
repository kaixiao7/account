package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
)

type CategorySrv interface {
	QueryAll(ctx context.Context, bookId int) ([]model.Category, error)
}

type categoryService struct {
	categoryStore store.CategoryStore
}

func NewCategorySrv() CategorySrv {
	return &categoryService{
		categoryStore: store.NewCategoryStore(),
	}
}

// QueryAll 查询账本下的所有分类
func (c *categoryService) QueryAll(ctx context.Context, bookId int) ([]model.Category, error) {
	return c.categoryStore.QueryAll(ctx, bookId)
}
