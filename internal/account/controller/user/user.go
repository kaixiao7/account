package user

import (
	"kaixiao7/account/internal/account/service"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/token"
)

type UserController struct {
	userSrv service.UserSrv
}

// NewUserController 创建用户处理器
func NewUserController() *UserController {
	return &UserController{userSrv: service.NewUserSrv()}
}

func generateTokens(userId int64) (*Tokens, error) {

	accessToken, err := token.GenerateAccessToken(userId)
	if err != nil {
		return nil, errno.NewWithError(errno.ErrToken, err)
	}
	refreshToken, err := token.GenerateRefreshToken(userId)
	if err != nil {
		return nil, errno.NewWithError(errno.ErrToken, err)
	}

	log.Debugf("access_token: %s \n refresh_token: %s\n", accessToken, refreshToken)
	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &tokens, nil
}
