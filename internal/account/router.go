package account

import (
	"kaixiao7/account/internal/account/controller/bill"
	"kaixiao7/account/internal/account/controller/book"
	"kaixiao7/account/internal/account/controller/budget"
	"kaixiao7/account/internal/account/controller/category"
	"kaixiao7/account/internal/account/controller/user"
	"kaixiao7/account/internal/pkg/core"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
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

	g.Use(middleware.Auth())

	users := g.Group("/users")
	{
		users.GET("/info", userController.Get)
	}

	books := g.Group("/books")
	{
		bookController := book.NewBookContorller()

		books.GET("", bookController.List)
	}

	budgets := g.Group("/books/:bookId/budgets")
	{
		budgetController := budget.NewBudgetController()

		budgets.GET("", budgetController.Get)
		budgets.PUT(":budgetId", budgetController.Put)
	}

	categories := g.Group("/books/:bookId/categories")
	{
		categoryController := category.NewCategoryController()

		categories.GET("", categoryController.List)
	}

	bills := g.Group("/books/:bookId/bills")
	{
		billController := bill.NewBillController()

		bills.POST("", billController.Add)
		bills.GET("", billController.List)
		bills.PUT(":billId", billController.Update)
		bills.DELETE(":billId", billController.Delete)
	}
}
