package model

// BorrowLendTotal 借入借出总计
type BorrowLendTotal struct {
	Borrow float64 `json:"borrow"`
	Lend   float64 `json:"lend"`
}

type BorrowLendFlow struct {
	Cost       float64 `json:"cost"`
	Type       int     `json:"type"`
	RecordTime int64   `json:"record_time"`
	// 账户id
	AccountId int64 `json:"account_id"`
	// 关联的借入借出id
	BorrowId int64  `json:"borrow_id"`
	Remark   string `json:"remark"`
}
