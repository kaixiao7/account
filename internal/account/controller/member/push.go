package member

import (
	"time"

	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (m *MemberController) Push(c *gin.Context) {

	var members []*model.Member
	if err := c.ShouldBindJSON(&members); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	syncTime := time.Now().Unix()

	if err := m.memberSrv.Push(c, members, syncTime); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, controller.PushRes{
		SyncTime: syncTime,
	})
}
