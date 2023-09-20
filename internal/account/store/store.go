package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kaixiao7/account/internal/pkg/constant"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Queryer interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)

	Get(dest any, query string, args ...any) error
	Select(dest any, query string, args ...any) error
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	QueryRowx(query string, args ...any) *sqlx.Row
	Preparex(query string) (*sqlx.Stmt, error)

	Rebind(query string) string
}

var db *sqlx.DB

// DBOption 数据库连接参数
type DBOption struct {
	Name                  string
	Host                  string
	Username              string
	Password              string
	Database              string
	Tls                   bool
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime int
}

func Init(option *DBOption) (*sqlx.DB, error) {
	var sqlxDB *sqlx.DB
	var err error
	if option.Name == "mysql" {
		sqlxDB, err = newMysqlDB(option)
	} else if option.Name == "pg" {
		sqlxDB, err = newPgDB(option)
	}

	if err != nil {
		return nil, err
	}

	db = sqlxDB

	return sqlxDB, nil
}

func newPgDB(opts *DBOption) (*sqlx.DB, error) {
	// postgres://username:password@localhost:5432/database_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opts.MaxOpenConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(opts.MaxConnectionLifeTime))

	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, nil
}

func newMysqlDB(opts *DBOption) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&tls=%t",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		opts.Tls)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(opts.MaxOpenConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(opts.MaxConnectionLifeTime))

	sqlxDB := sqlx.NewDb(db, "mysql")
	return sqlxDB, nil
}

// func newSqliteDB(opts *DBOption) (*sqlx.DB, error) {
// 	dsn := fmt.Sprintf("file:%s", opts.File)
// 	sqlDB, err := sql.Open("sqlite3", dsn)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
// 	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
// 	sqlDB.SetConnMaxLifetime(time.Duration(opts.MaxConnectionLifeTime))
//
// 	sqlxDB := sqlx.NewDb(sqlDB, "sqlite3")
//
// 	return sqlxDB, nil
// }

func Close() error {
	if db == nil {
		return nil
	}

	return db.Close()
}

func getDBFromContext(ctx context.Context) Queryer {
	db := ctx.Value(constant.SqlDBKey)
	return db.(Queryer)
}

var _ Queryer = (*sqlx.DB)(nil)
var _ Queryer = (*sqlx.Tx)(nil)
