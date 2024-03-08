package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
)

type CategoryTagSrv interface {
	Push(ctx context.Context, tags []*model.CategoryTag, syncTime int64) error
	Pull(ctx context.Context, bookId int64, lastSyncTime int64) ([]*model.CategoryTag, error)
}

type categoryTagService struct {
	tagStore store.CategoryTagStore
}

func NewCategoryTagSrv() CategoryTagSrv {
	return &categoryTagService{
		tagStore: store.NewCategoryTagStore(),
	}
}

func (c *categoryTagService) Push(ctx context.Context, tags []*model.CategoryTag, syncTime int64) error {
	return WithTransaction(ctx, func(ctx context.Context) error {
		for _, tag := range tags {
			tag.SyncTime = syncTime
			if tag.SyncState == constant.SYNC_ADD {
				tag.SyncState = constant.SYNC_SUCCESS
				if e := c.tagStore.Add(ctx, tag); e != nil {
					return e
				}
			} else {
				tag.SyncState = constant.SYNC_SUCCESS
				if e := c.tagStore.Update(ctx, tag); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (c *categoryTagService) Pull(ctx context.Context, bookId int64, lastSyncTime int64) ([]*model.CategoryTag, error) {
	return c.tagStore.QueryBySyncTime(ctx, bookId, lastSyncTime)
}
