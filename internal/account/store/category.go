package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type CategoryStore interface {
	Add(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Category, error)

	QueryAll(ctx context.Context, bookId int64) ([]model.Category, error)
	// QueryById 通过主键查询
	QueryById(ctx context.Context, id int64) (*model.Category, error)

	// QueryByUserId 根据用户id查询其所有的分类
	QueryByUserId(ctx context.Context, userId int64) ([]model.Category, error)
}

type category struct {
}

func NewCategoryStore() CategoryStore {
	return &category{}
}

// QueryAll 查询账本下的所有分类
func (c *category) QueryAll(ctx context.Context, bookId int64) ([]model.Category, error) {
	db := getDBFromContext(ctx)

	sql := db.Rebind("select * from category where book_id = ? order by sort")
	var categories = []model.Category{}
	err := db.Select(&categories, sql, bookId)
	if err != nil {
		return nil, errors.Wrap(err, "query category all store.")
	}

	return categories, nil
}

// QueryByUserId 根据用户id查询其所有的分类
func (c *category) QueryByUserId(ctx context.Context, userId int64) ([]model.Category, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from category where user_id = ?")
	var categories = []model.Category{}
	err := db.Select(&categories, querySql, userId)

	if err != nil {
		return nil, errors.Wrap(err, "query category by userId store.")
	}

	return categories, nil
}

// QueryById 通过主键查询
func (c *category) QueryById(ctx context.Context, id int64) (*model.Category, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from category where id = ?")
	var category model.Category
	err := db.Get(&category, querySql, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query category by id store.")
	}

	return &category, nil
}

func (c *category) Add(ctx context.Context, category *model.Category) error {
	db := getDBFromContext(ctx)

	sql := db.Rebind(`insert into category(id, name, icon, color, sort, type, book_id, user_id, del_flag, sync_state, sync_time,
				create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(sql, category.Id, category.Name, category.Icon, category.Color, category.Sort, category.Type,
		category.BookId, category.UserId, category.DelFlag, category.SyncState, category.SyncTime, category.CreateTime, category.UpdateTime)

	if err != nil {
		return errors.Wrap(err, "add category store.")
	}

	return nil
}

// Update 更新分类
func (c *category) Update(ctx context.Context, category *model.Category) error {
	db := getDBFromContext(ctx)

	sql := db.Rebind(`update category set name=?,icon=?,color=?,sort=?,type=?,book_id=?,user_id=?,del_flag=?,sync_state=?,sync_time=?,update_time=? where id=?`)
	_, err := db.Exec(sql, category.Name, category.Icon, category.Color, category.Sort, category.Type, category.BookId,
		category.UserId, category.DelFlag, category.SyncState, category.SyncTime, category.UpdateTime, category.Id)
	if err != nil {
		return errors.Wrap(err, "update category store.")
	}
	return nil
}

func (c *category) QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Category, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from category where book_id = ? and sync_time > ?")
	var category = []*model.Category{}
	err := db.Select(&category, querySql, bookId, syncTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query category sync time store.")
	}

	return category, nil
}
