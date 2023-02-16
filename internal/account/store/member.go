package store

import (
	"context"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type MemberStore interface {
	Add(ctx context.Context, member *model.Member) error

	Update(ctx context.Context, member *model.Member) error

	QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Member, error)

	QueryByBookId(ctx context.Context, bookId int64) ([]*model.Member, error)
}

type member struct {
}

func NewMemberStore() MemberStore {
	return &member{}
}

func (m *member) Add(ctx context.Context, member *model.Member) error {
	db := getDBFromContext(ctx)

	insertSql := "insert into book_member values(?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(insertSql, member.Id, member.BookId, member.UserId, member.Username, member.DelFlag, member.SyncState, member.SyncTime,
		member.CreateTime, member.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "member add store")
	}

	return nil
}

func (m *member) Update(ctx context.Context, member *model.Member) error {
	db := getDBFromContext(ctx)

	updateSql := "update book_member set book_id=?, user_id=?, username=?, del_flag=?, sync_state=?, sync_time=?, create_time=?,update_time=? where id=?"
	_, err := db.Exec(updateSql, member.BookId, member.UserId, member.Username, member.DelFlag, member.SyncState, member.SyncTime,
		member.CreateTime, member.UpdateTime, member.Id)

	if err != nil {
		return errors.Wrap(err, "member update store")
	}
	return nil
}

func (m *member) QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Member, error) {
	db := getDBFromContext(ctx)

	querySql := `select * from book_member where book_id = ? and sync_time > ?`

	var memberList = []*model.Member{}
	err := db.Select(&memberList, querySql, bookId, syncTime)

	if err != nil {
		return nil, errors.Wrap(err, "query member sync time store")
	}

	return memberList, nil
}

func (m *member) QueryByBookId(ctx context.Context, bookId int64) ([]*model.Member, error) {
	db := getDBFromContext(ctx)

	querySql := `select * from book_member where book_id = ?`

	var memberList = []*model.Member{}
	err := db.Select(&memberList, querySql, bookId)

	if err != nil {
		return nil, errors.Wrap(err, "query member list store")
	}

	return memberList, nil
}
