package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BookStore interface {
	Add(ctx context.Context, book *model.Book) error

	Update(ctx context.Context, book *model.Book) error

	QueryBySyncTime(ctx context.Context, userId int, syncTime int64) ([]*model.Book, error)

	QueryBookList(ctx context.Context, userId int) ([]*model.Book, error)

	// QueryById 根据主键id查询
	QueryById(ctx context.Context, id int) (*model.Book, error)

	// QueryBookMember 根据账本id查询该账本下的所有成员
	// QueryBookMember(ctx context.Context, bookId int) ([]int, error)
}

type book struct {
}

func NewBookStore() BookStore {
	return &book{}
}

func (b *book) Add(ctx context.Context, book *model.Book) error {
	db := getDBFromContext(ctx)

	insertSql := "insert into user_book values(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(insertSql, book.Id, book.BookName, book.UserId, book.Cover, book.DelFlag, book.SyncState, book.SyncTime,
		book.CreateTime, book.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "book add store")
	}

	return nil
}

func (b *book) Update(ctx context.Context, book *model.Book) error {
	db := getDBFromContext(ctx)

	updateSql := "update user_book set book_name=?, user_id=?, cover=?, del_flag=?, sync_state=?, sync_time=?, create_time=?,update_time=? where id=?"
	_, err := db.Exec(updateSql, book.BookName, book.UserId, book.Cover, book.DelFlag, book.SyncState, book.SyncTime,
		book.CreateTime, book.UpdateTime, book.Id)

	if err != nil {
		return errors.Wrap(err, "book update store")
	}
	return nil
}

func (b *book) QueryBySyncTime(ctx context.Context, userId int, syncTime int64) ([]*model.Book, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from user_book where user_id = ? and sync_time > ?"

	var bookList = []*model.Book{}
	err := db.Select(&bookList, querySql, userId, syncTime)

	if err != nil {
		return nil, errors.Wrap(err, "query book sync time store")
	}

	return bookList, nil
}

// QueryBookList 查询用户账本列表
func (b *book) QueryBookList(ctx context.Context, userId int) ([]*model.Book, error) {
	db := getDBFromContext(ctx)
	sql := `
		select *
		from user_book
		where del_flag = 0
		and id in (
			select book_id
			from book_member
			where user_id = ?
		)
	`
	var bookList = []*model.Book{}
	err := db.Select(&bookList, sql, userId)

	if err != nil {
		return nil, errors.Wrap(err, "query book list store")
	}

	return bookList, nil
}

// QueryById 根据主键id查询
func (b *book) QueryById(ctx context.Context, id int) (*model.Book, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from user_book where id = ?"
	var book model.Book
	err := db.Get(&book, querySql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query book by id store")
	}

	return &book, nil
}

// QueryBookMember 根据账本id查询该账本下的所有成员
// func (b *book) QueryBookMember(ctx context.Context, bookId int) ([]int, error) {
// 	db := getDBFromContext(ctx)
//
// 	querySql := "select user_id from book_member where book_id = ?"
// 	var memberIds []int
// 	err := db.Select(&memberIds, querySql, bookId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, errors.Wrap(err, "query book member store")
// 	}
//
// 	return memberIds, nil
// }
