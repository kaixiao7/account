package model

type AssetFlow struct {
	Id            int     `db:"id" json:"id,omitempty"`
	UserId        int     `db:"user_id" json:"user_id,omitempty"`
	AssetId       int     `db:"asset_id" json:"asset_id"`
	Type          int     `db:"type" json:"type"`
	Cost          float64 `db:"cost" json:"cost"`
	RecordTime    int64   `db:"record_time" json:"record_time"`
	Del           int     `db:"del" json:"-"`
	Remark        string  `db:"remark" json:"remark,omitempty"`
	TargetAssetId *int    `db:"target_asset_id" json:"target_asset_id,omitempty"`
	AssociateName string  `db:"associate_name" json:"associate_name,omitempty"`
	Finished      *int    `db:"finished" json:"finished,omitempty"`
	BorrowId      *int    `db:"borrow_id" json:"borrow_id,omitempty"`
	CreateTime    int64   `db:"create_time" json:"create_time,omitempty"`
	UpdateTime    int64   `db:"update_time" json:"update_time,omitempty"`
}

type AssetFlowVO struct {
	Id            int     `json:"id,omitempty"`
	UserId        int     `json:"user_id,omitempty"`
	AssetId       int     `json:"asset_id"`
	Type          int     `json:"type"`
	Cost          float64 `json:"cost"`
	RecordTime    int64   `json:"record_time"`
	Remark        string  `json:"remark,omitempty"`
	CategoryId    int     `json:"category_id,omitempty"`
	TargetAssetId int     `json:"target_asset_id,omitempty"`
	AssociateName string  `json:"associate_name,omitempty"`
	CreateTime    int64   `json:"create_time,omitempty"`
	UpdateTime    int64   `json:"update_time,omitempty"`
}

func AssetFlow2VO(flows []AssetFlow) []AssetFlowVO {
	var ret = []AssetFlowVO{}
	for _, flow := range flows {
		ret = append(ret, AssetFlowVO{
			Id:            flow.Id,
			UserId:        flow.UserId,
			AssetId:       flow.AssetId,
			Type:          flow.Type,
			Cost:          flow.Cost,
			RecordTime:    flow.RecordTime,
			Remark:        flow.Remark,
			TargetAssetId: *flow.TargetAssetId,
			AssociateName: flow.AssociateName,
			CreateTime:    flow.CreateTime,
			UpdateTime:    flow.UpdateTime,
		})
	}
	return ret
}

func Bill2VO(bills []Bill) []AssetFlowVO {
	var ret = []AssetFlowVO{}
	for _, bill := range bills {
		ret = append(ret, AssetFlowVO{
			Id:         bill.Id,
			UserId:     bill.UserId,
			AssetId:    bill.AssetId,
			Type:       int(*bill.Type),
			Cost:       bill.Cost,
			RecordTime: bill.RecordTime,
			Remark:     bill.Remark,
			CategoryId: bill.CategoryId,
			CreateTime: bill.CreateTime,
			UpdateTime: bill.UpdateTime,
		})
	}
	return ret
}
