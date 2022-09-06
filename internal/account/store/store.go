package store

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

// DBOption 数据库连接参数
type DBOption struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime int
}

func Init(option *DBOption) error {
	g, err := newDB(option)
	if err != nil {
		return err
	}

	db = g
	return nil
}

func newDB(opts *DBOption) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDb.SetMaxIdleConns(opts.MaxIdleConnections)
	sqlDb.SetConnMaxLifetime(time.Duration(opts.MaxConnectionLifeTime))

	return db, nil
}

func Close() error {
	if db == nil {
		return nil
	}

	d, err := db.DB()
	if err != nil {
		return err
	}

	return d.Close()
}
