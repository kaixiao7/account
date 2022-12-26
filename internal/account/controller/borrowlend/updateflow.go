package borrowlend

import (
	"kaixiao7/account/internal/pkg/timex"

	"github.com/gin-gonic/gin"
)

type updateFlowReq struct {
	AccountId  int            `json:"account_id" binding:"required,numeric"`
	Cost       float64        `json:"cost" binding:"required,numeric"`
	RecordTime timex.JsonTime `json:"record_time" binding:"required"`
	Type       int            `json:"type,omitempty" binding:"required,numeric"`
	Remark     string         `json:"remark,omitempty"`
}

func (b *BorrowLendController) UpdateFlow(c *gin.Context) {

}
