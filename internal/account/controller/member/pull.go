package member

import (
	"kaixiao7/account/internal/account/controller"
	"kaixiao7/account/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

func (m *MemberController) Pull(c *gin.Context) {
	bookId, exist := controller.GetInt64ParamFromUrl(c, "bookId")
	if !exist {
		return
	}

	lastSyncTime, exist := controller.GetInt64ParamFromParam(c, "lastSyncTime")
	if !exist {
		return
	}

	members, err := m.memberSrv.Pull(c, bookId, lastSyncTime)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, members)
}
