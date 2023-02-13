package model

type Category struct {
	Id         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Icon       string `db:"icon" json:"icon"`
	Color      string `db:"color" json:"color"`
	Sort       int    `db:"sort" json:"sort"`
	Type       int    `db:"type" json:"type"`
	BookId     int    `db:"book_id" json:"book_id"`
	UserId     int    `db:"user_id" json:"user_id"`
	DelFlag    int    `db:"del_flag" json:"del_flag"`
	SyncState  int    `db:"sync_state" json:"sync_state"`
	SyncTime   int64  `db:"sync_time" json:"sync_time"`
	CreateTime int64  `db:"create_time" json:"create_time"`
	UpdateTime int64  `db:"update_time" json:"update_time"`
}
