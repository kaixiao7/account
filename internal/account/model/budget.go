package model

type Budget struct {
	Id         int     `db:"id" json:"id"`
	Budget     float64 `db:"budget" json:"budget"`
	BookId     int     `db:"book_id" json:"book_id"`
	Type       int     `db:"type" json:"type"`
	CategoryId int     `db:"category_id" json:"category_id"`
	CreateId   int     `db:"create_id" json:"create_id"`
	SyncState  int     `db:"sync_state" json:"sync_state"`
	SyncTime   int64   `db:"sync_time" json:"sync_time"`
	CreateTime int64   `db:"create_time" json:"create_time"`
	UpdateTime int64   `db:"update_time" json:"update_time"`
}
