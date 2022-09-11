package bill

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

type billReq struct {
	Cost       float32        `db:"cost" json:"cost,omitempty" binding:"required,numeric"`
	Type       *int8          `db:"type" json:"type,omitempty" binding:"required,gte=0,lte=1"`
	Remark     string         `db:"remark" json:"remark,omitempty" binding:"required,max=200"`
	RecordTime timex.JsonTime `db:"record_time" json:"record_time" binding:"required"`
	AccountId  int            `db:"account_id" json:"account_id,omitempty" binding:"required,numeric"`
	CategoryId int            `db:"category_id" json:"category_id,omitempty" binding:"required,min=1"`
}

func (b *BillController) Add(c *gin.Context) {
	userId := controller.GetUserId(c)
	bookId, ok := controller.GetIntParamFromUrl(c, "bookId")
	if !ok {
		return
	}

	var req billReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	bill := model.Bill{
		Cost:       req.Cost,
		Type:       req.Type,
		Remark:     req.Remark,
		RecordTime: req.RecordTime.Timestamp(),
		AccountId:  req.AccountId,
		CategoryId: req.CategoryId,
		BookId:     bookId,
		UserId:     userId,
	}
	err := b.billSrv.Add(c, &bill)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
