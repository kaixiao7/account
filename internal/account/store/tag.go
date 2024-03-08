package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type CategoryTagStore interface {
	Add(ctx context.Context, tag *model.CategoryTag) error
	Update(ctx context.Context, tag *model.CategoryTag) error

	QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.CategoryTag, error)
}

type categoryTag struct {
}

func NewCategoryTagStore() CategoryTagStore {
	return &categoryTag{}
}

func (c *categoryTag) Add(ctx context.Context, tag *model.CategoryTag) error {
	db := getDBFromContext(ctx)

	sql := db.Rebind(`insert into category_tag(id, category_id, book_id, user_id, name, weight, sync_state, sync_time,
				create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(sql, tag.Id, tag.CategoryId, tag.BookId, tag.UserId, tag.Name, tag.Weight, tag.SyncState,
		tag.SyncTime, tag.CreateTime, tag.UpdateTime)

	if err != nil {
		return errors.Wrap(err, "add category tag store.")
	}

	return nil
}

// Update 更新分类
func (c *categoryTag) Update(ctx context.Context, tag *model.CategoryTag) error {
	db := getDBFromContext(ctx)

	sql := db.Rebind(`update category_tag set category_id=?,book_id=?,user_id=?,name=?,weight=?,sync_state=?,
                        sync_time=?,update_time=? where id=?`)
	_, err := db.Exec(sql, tag.CategoryId, tag.BookId, tag.UserId, tag.Name, tag.Weight, tag.SyncState, tag.SyncTime,
		tag.UpdateTime, tag.Id)
	if err != nil {
		return errors.Wrap(err, "update category tag store.")
	}
	return nil
}

func (c *categoryTag) QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.CategoryTag, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from category_tag where book_id = ? and sync_time > ?")
	var tag = []*model.CategoryTag{}
	err := db.Select(&tag, querySql, bookId, syncTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query category tag sync time store.")
	}

	return tag, nil
}
