package user

import (
	"github.com/gin-gonic/gin"
	"kaixiao7/account/internal/pkg/auth"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/token"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// Login 用户登录
func (u *UserController) Login(c *gin.Context) {
	var r LoginRequest

	// 绑定数据
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteRespErr(c, err)
		return
	}

	// 数据库查询用户信息
	user, err := u.userSrv.Get(r.Username)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	// 比较密码
	if err := auth.Compare(user.Password, r.Password); err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrPasswordIncorrect))
		return
	}

	// 登录成功，发送token
	t, err := token.Sign(user.Id)
	if err != nil {
		core.WriteRespErr(c, errno.NewWithError(errno.ErrToken, err))
		return
	}

	core.WriteRespSuccess(c, LoginResponse{Token: t})
}
