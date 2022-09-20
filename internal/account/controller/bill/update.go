package bill

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (b *BillController) Update(c *gin.Context) {
	userId := controller.GetUserId(c)

	bookId, ok := controller.GetIntParamFromUrl(c, "bookId")
	if !ok {
		return
	}

	billId, ok := controller.GetIntParamFromUrl(c, "billId")
	if !ok {
		return
	}

	var req billReq
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	bill := model.Bill{
		Id:         billId,
		Cost:       req.Cost,
		Type:       req.Type,
		Remark:     req.Remark,
		RecordTime: req.RecordTime.Timestamp(),
		AssetId:    req.AssetId,
		CategoryId: req.CategoryId,
		BookId:     bookId,
		UserId:     userId,
	}

	if err := b.billSrv.Update(c, &bill); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, nil)
}
