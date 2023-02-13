package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
)

type CategorySrv interface {
	Push(ctx context.Context, categories []*model.Category, syncTime int64) error
	Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Category, error)

	QueryAll(ctx context.Context, bookId int) ([]model.Category, error)

	// QueryByUserId 根据用户id查询其所有分类
	QueryByUserId(ctx context.Context, userId int) ([]model.Category, error)
}

type categoryService struct {
	categoryStore store.CategoryStore
}

func NewCategorySrv() CategorySrv {
	return &categoryService{
		categoryStore: store.NewCategoryStore(),
	}
}

func (c *categoryService) Push(ctx context.Context, categories []*model.Category, syncTime int64) error {
	return WithTransaction(ctx, func(ctx context.Context) error {
		for _, category := range categories {
			category.SyncTime = syncTime
			if category.SyncState == constant.SYNC_ADD {
				category.SyncState = constant.SYNC_SUCCESS
				if e := c.categoryStore.Add(ctx, category); e != nil {
					return e
				}
			} else {
				category.SyncState = constant.SYNC_SUCCESS
				if e := c.categoryStore.Update(ctx, category); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (c *categoryService) Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Category, error) {
	return c.categoryStore.QueryBySyncTime(ctx, bookId, lastSyncTime)
}

// QueryAll 查询账本下的所有分类
func (c *categoryService) QueryAll(ctx context.Context, bookId int) ([]model.Category, error) {
	return c.categoryStore.QueryAll(ctx, bookId)
}

// QueryByUserId 根据用户id查询其所有分类
func (c *categoryService) QueryByUserId(ctx context.Context, userId int) ([]model.Category, error) {
	return c.categoryStore.QueryByUserId(ctx, userId)
}
