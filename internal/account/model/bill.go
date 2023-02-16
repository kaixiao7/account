package model

type BillTag struct {
	CategoryId int64  `db:"category_id" json:"category_id"`
	Remark     string `db:"remark" json:"tag"`
}
