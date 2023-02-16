package model

type Member struct {
	Id         int64  `db:"id" json:"id"`
	BookId     int64  `db:"book_id" json:"book_id"`
	UserId     int64  `db:"user_id" json:"user_id"`
	Username   string `db:"username" json:"username"`
	DelFlag    int    `db:"del_flag" json:"del_flag"`
	SyncState  int    `db:"sync_state" json:"sync_state"`
	SyncTime   int64  `db:"sync_time" json:"sync_time"`
	CreateTime int64  `db:"create_time" json:"create_time"`
	UpdateTime int64  `db:"update_time" json:"update_time"`
}
