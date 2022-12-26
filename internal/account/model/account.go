package model

type Account struct {
	Id          int     `db:"id" json:"id,omitempty"`
	UserId      int     `db:"user_id" json:"user_id,omitempty"`
	AccountType *int    `db:"account_type" json:"account_type" binding:"required,numeric"`
	AccountName string  `db:"account_name" json:"account_name" binding:"required,max=30"`
	Balance     float64 `db:"balance" json:"balance" binding:"required,numeric"`
	Init        float64 `db:"init" json:"-"`
	Icon        string  `db:"icon" json:"icon" binding:"required"`
	Del         int     `db:"del" json:"-"`
	IsTotal     int     `db:"is_total" json:"is_total" binding:"required"`
	Remark      string  `db:"remark" json:"remark,omitempty" binding:"max=200"`
	CreateTime  int64   `db:"create_time" json:"create_time,omitempty"`
	UpdateTime  int64   `db:"update_time" json:"update_time,omitempty"`
}
