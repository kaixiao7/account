package service

import (
	"context"

	"kaixiao7/account/internal/pkg/constant"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type handler func(ctx context.Context) error

// WithTransaction 在事务环境下执行相关操作
func WithTransaction(ctx context.Context, h handler) error {
	db := ctx.Value(constant.SqlDBKey).(*sqlx.DB)

	tx, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "get transaction fail.")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	c := context.WithValue(ctx, constant.SqlDBKey, tx)

	if err := h(c); err != nil {
		return errors.Wrap(err, "transaction exec fail.")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "transaction commit fail.")
	}
	return nil
}
