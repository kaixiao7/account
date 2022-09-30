package user

import (
	"kaixiao7/account/internal/account/service"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/token"

	"github.com/spf13/viper"
)

type UserController struct {
	userSrv service.UserSrv
}

// NewUserController 创建用户处理器
func NewUserController() *UserController {
	return &UserController{userSrv: service.NewUserSrv()}
}

func generateTokens(userId int) (*Tokens, error) {
	// token
	tokenExp := viper.GetInt("jwt.expire")
	t, err := token.Sign(userId, tokenExp)
	if err != nil {
		return nil, errno.NewWithError(errno.ErrToken, err)
	}

	// refreshToken
	refreshTokenExp := viper.GetInt("jwt.refresh-token")
	refreshToken, err := token.Sign(userId, refreshTokenExp*24)
	if err != nil {
		return nil, errno.NewWithError(errno.ErrToken, err)
	}

	log.Debugf("access_token: %s \n refresh_token: %s\n", t, refreshToken)
	tokens := Tokens{
		AccessToken:  t,
		RefreshToken: refreshToken,
	}

	return &tokens, nil
}
