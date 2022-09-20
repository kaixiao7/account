package model

type Bill struct {
	Id         int     `db:"id" json:"id,omitempty"`
	Cost       float32 `db:"cost" json:"cost,omitempty" binding:"required,numeric"`
	Type       *int8   `db:"type" json:"type,omitempty" binding:"required,gte=0,lte=1"`
	Remark     string  `db:"remark" json:"remark,omitempty" binding:"required,max=200"`
	RecordTime int64   `db:"record_time" json:"record_time" binding:"required"`
	UserId     int     `db:"user_id" json:"user_id,omitempty"`
	BookId     int     `db:"book_id" json:"book_id,omitempty" binding:"required,gte=1"`
	AssetId    int     `db:"asset_id" json:"asset_id,omitempty" binding:"required,numeric"`
	CategoryId int     `db:"category_id" json:"category_id,omitempty" binding:"required,min=1"`
	Username   string  `db:"user_name" json:"username,omitempty" binding:"min=alphanum,max=20"`
	CreateTime int64   `db:"create_time" json:"create_time,omitempty"`
	UpdateTime int64   `db:"update_time" json:"update_time,omitempty"`
}

type BillTag struct {
	CategoryId int    `db:"category_id" json:"category_id"`
	Remark     string `db:"remark" json:"tag"`
}
