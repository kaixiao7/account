package model

type BillTag struct {
	CategoryId int    `db:"category_id" json:"category_id"`
	Remark     string `db:"remark" json:"tag"`
}
