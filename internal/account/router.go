package account

import (
	"github.com/gin-gonic/gin"
	"kaixiao7/account/internal/account/controller/user"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/middleware"
)

func loadRouter(g *gin.Engine, mw ...gin.HandlerFunc) {
	installMiddleware(g, mw...)
	installController(g)
}

func installMiddleware(g *gin.Engine, mw ...gin.HandlerFunc) {
	g.Use(gin.Recovery(), gin.Logger())
	g.Use(mw...)
}

func installController(g *gin.Engine) {
	// 404
	g.NoRoute(func(c *gin.Context) {
		core.WriteRespErr(c, errno.New(errno.ErrPageNotFound))
	})

	userController := user.NewUserController()

	g.POST("/login", userController.Login)

	users := g.Group("/user", middleware.Auth())
	{
		users.GET("/info", userController.Get)
	}
}
