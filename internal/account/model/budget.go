package model

type Budget struct {
	Id         int
	Budget     float32
	BookId     int `db:"book_id" json:"book_id"`
	Type       int
	CategoryId int   `db:"category_id" json:"category_id"`
	CreateId   int   `db:"create_id" json:"create_id"`
	CreateTime int64 `db:"create_time" json:"create_time"`
	UpdateTime int64 `db:"update_time" json:"update_time"`
}
