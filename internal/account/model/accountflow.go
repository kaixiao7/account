package model

type AccountFlow struct {
	Id              int      `db:"id" json:"id,omitempty"`
	UserId          int      `db:"user_id" json:"user_id,omitempty"`
	Username        string   `db:"username" json:"username"`
	AccountId       int      `db:"account_id" json:"account_id"`
	Type            int      `db:"type" json:"type"`
	Cost            float64  `db:"cost" json:"cost"`
	RecordTime      int64    `db:"record_time" json:"record_time"`
	DelFlag         int      `db:"del_flag" json:"del_flag"`
	BookId          *int     `db:"book_id" json:"book_id,omitempty"`
	CategoryId      *int     `db:"category_id" json:"category_id,omitempty"`
	Remark          string   `db:"remark" json:"remark,omitempty"`
	TargetAccountId *int     `db:"target_account_id" json:"target_account_id,omitempty"`
	AssociateName   string   `db:"associate_name" json:"associate_name,omitempty"`
	Finished        *int     `db:"finished" json:"finished,omitempty"`
	BorrowLendId    *int     `db:"borrow_lend_id" json:"borrow_lend_id,omitempty"`
	Profit          *float64 `db:"profit" json:"profit,omitempty"`
	SyncState       int      `db:"sync_state" json:"sync_state"`
	SyncTime        int64    `db:"sync_time" json:"sync_time"`
	CreateTime      int64    `db:"create_time" json:"create_time,omitempty"`
	UpdateTime      int64    `db:"update_time" json:"update_time,omitempty"`
}
