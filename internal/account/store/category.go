package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type CategoryStore interface {
	QueryAll(ctx context.Context, bookId int) ([]model.Category, error)
	Add(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error

	// QueryById 通过主键查询
	QueryById(ctx context.Context, id int) (*model.Category, error)

	// QueryByUserId 根据用户id查询其所有的分类
	QueryByUserId(ctx context.Context, userId int) ([]model.Category, error)
}

type category struct {
}

func NewCategoryStore() CategoryStore {
	return &category{}
}

// QueryAll 查询账本下的所有分类
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

// QueryByUserId 根据用户id查询其所有的分类
func (c *category) QueryByUserId(ctx context.Context, userId int) ([]model.Category, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from category where user_id = ?"
	var categories = []model.Category{}
	err := db.Select(&categories, querySql, userId)

	if err != nil {
		return nil, errors.Wrap(err, "query category by userId store.")
	}

	return categories, nil
}

// QueryById 通过主键查询
func (c *category) QueryById(ctx context.Context, id int) (*model.Category, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from category where id = ?"
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
