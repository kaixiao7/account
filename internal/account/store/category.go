package store

import (
	"context"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type CategoryStore interface {
	QueryAll(ctx context.Context, bookId int) ([]model.Category, error)
	Add(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
}

type category struct {
}

func NewCategoryStore() CategoryStore {
	return &category{}
}

// QueryAll 查询所有分类
func (c *category) QueryAll(ctx context.Context, bookId int) ([]model.Category, error) {
	db := getDBFromContext(ctx)

	sql := "select * from category where book_id = ? order by sort"
	var categories = []model.Category{}
	err := db.Select(&categories, sql, bookId)
	if err != nil {
		return nil, errors.Wrap(err, "query category all store.")
	}

	return categories, nil
}

func (c *category) Add(ctx context.Context, category *model.Category) error {
	db := getDBFromContext(ctx)

	sql := "insert into category(name, icon, color, sort, type, book_id, user_id, create_time, update_time)" +
		" values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sql, category.Name, category.Icon, category.Color, category.Sort, category.Type,
		category.BookId, category.UserId, category.CreateTime, category.UpdateTime)

	if err != nil {
		return errors.Wrap(err, "add category store.")
	}

	return nil
}

// Update 更新分类
func (c *category) Update(ctx context.Context, category *model.Category) error {
	db := getDBFromContext(ctx)

	sql := "update category set name=?,icon=?,color=?,sort=?,user_id=?,update_time=? where id=?"
	_, err := db.Exec(sql, category.Name, category.Color, category.Sort, category.UserId, category.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "update category store.")
	}
	return nil
}
