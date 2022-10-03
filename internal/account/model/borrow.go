package model

// BorrowTotal 借入借出总计
type BorrowTotal struct {
	BorrowIn  float64 `json:"borrow_in"`
	BorrowOut float64 `json:"borrow_out"`
}

type BorrowFlow struct {
	Id          int     `db:"id" json:"id"`
	AssetFlowId int     `db:"asset_flow_id" json:"asset_flow_id"`
	AssetId     int     `db:"asset_id" json:"asset_id"`
	Cost        float64 `db:"cost" json:"cost"`
	RecordTime  int64   `db:"record_time" json:"record_time"`
	Type        int     `db:"type" json:"type"`
	Remark      string  `db:"remark" json:"remark"`
	CreateTime  int64   `db:"create_time" json:"create_time"`
	UpdateTime  int64   `db:"update_time" json:"update_time"`
}
