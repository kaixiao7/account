package model

type Book struct {
	Id         int    `db:"id" json:"id"`
	BookName   string `db:"book_name" json:"book_name"`
	UserId     int    `db:"user_id" json:"user_id"`
	Cover      string `db:"cover" json:"cover"`
	DelFlag    int    `db:"del_flag" json:"-"`
	CreateTime int64  `db:"create_time" json:"create_time"`
	UpdateTime int64  `db:"update_time" json:"update_time"`
}
