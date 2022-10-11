package user

import (
	"kaixiao7/account/internal/pkg/auth"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	user, err := u.userSrv.Get(c, r.Username)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	// 比较密码
	if err := auth.Compare(user.Password, r.Password); err != nil {
		core.WriteRespErr(c, errno.New(errno.ErrPasswordIncorrect))
		return
	}

	// 登录成功，生成token并响应
	resp, err := generateTokens(user.Id)
	if err != nil {
		core.WriteRespErr(c, err)
		return
	}

	core.WriteRespSuccess(c, resp)
}
