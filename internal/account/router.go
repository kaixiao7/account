package account

import (
	"kaixiao7/account/internal/account/controller/account"
	"kaixiao7/account/internal/account/controller/accountflow"
	"kaixiao7/account/internal/account/controller/book"
	"kaixiao7/account/internal/account/controller/budget"
	"kaixiao7/account/internal/account/controller/category"
	"kaixiao7/account/internal/account/controller/member"
	"kaixiao7/account/internal/account/controller/tag"
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

	g.GET("/refresh", userController.RefreshToken)

	users := g.Group("/users")
	{
		users.GET("/info", userController.Get)
		users.PUT("", userController.Update)
		users.PUT("/changePwd", userController.ChangePassword)
	}

	books := g.Group("/books")
	{
		bookController := book.NewBookContorller()

		books.GET("", bookController.List)
		books.GET("/pull", bookController.Pull)
		books.POST("/push", bookController.Push)
	}

	budgets := g.Group("/books/:bookId/budgets")
	{
		budgetController := budget.NewBudgetController()

		budgets.GET("", budgetController.Get)
		budgets.PUT(":budgetId", budgetController.Put)
		budgets.GET("/pull", budgetController.Pull)
		budgets.POST("/push", budgetController.Push)
	}

	members := g.Group("/books/:bookId/members")
	{
		memberController := member.NewMemberController()

		members.GET("/pull", memberController.Pull)
		members.POST("/push", memberController.Push)
	}

	categories := g.Group("/books/:bookId/categories")
	{
		categoryController := category.NewCategoryController()

		categories.GET("", categoryController.List)
		categories.GET("/pull", categoryController.Pull)
		categories.POST("/push", categoryController.Push)

		g.GET("/categories", categoryController.ListUser)
	}

	tags := g.Group("/books/:bookId/tags")
	{
		tagController := tag.NewCategoryTagController()

		tags.GET("/pull", tagController.Pull)
		tags.POST("/push", tagController.Push)
	}

	// bills := g.Group("/books/:bookId/bills")
	// {
	// 	billController := bill.NewBillController()
	//
	// 	bills.POST("", billController.Add)
	// 	bills.GET("", billController.List)
	// 	bills.PUT(":billId", billController.Update)
	// 	bills.DELETE(":billId", billController.Delete)
	// }
	//
	// tags := g.Group("/tags")
	// {
	// 	tagController := tag.NewTagController()
	//
	// 	// 查询标签，需要查询参数bookId，/tags?bookId=1
	// 	tags.GET("", tagController.List)
	// }

	accounts := g.Group("/accounts")
	{
		accountController := account.NewAccountController()

		accounts.GET("", accountController.List)
		accounts.POST("", accountController.Add)
		accounts.PUT(":accountId", accountController.Update)
		accounts.DELETE(":accountId", accountController.Delete)
		accounts.GET("/pull", accountController.Pull)
		accounts.POST("/push", accountController.Push)
	}

	accountFlows := g.Group("/accounts/:accountId/flows")
	{
		accountFlowController := accountflow.NewAccountFlowController()

		accountFlows.GET("", accountFlowController.List)
		accountFlows.POST("", accountFlowController.Add)
		// accountFlows.PUT(":accountFlowId", accountFlowController.Update)
		accountFlows.DELETE(":accountFlowId", accountFlowController.Delete)

		g.GET("/flows/:bookId/pull", accountFlowController.PullBook)
		g.GET("/flows/pull", accountFlowController.PullOther)
		g.POST("/flows/push", accountFlowController.Push)
	}

	// borrows := g.Group("/borrowlends")
	// {
	// 	borrowController := borrowlend.NewBorrowController()
	//
	// 	borrows.GET("", borrowController.List)
	// 	borrows.GET("/total", borrowController.Total)
	//
	// 	borrows.GET("/:accountFlowId/flows", borrowController.FlowList)
	// 	borrows.POST("/:accountFlowId/flows", borrowController.AddFlow)
	// 	// borrows.PUT("/:accountFlowId/flows/:flowId", borrowController.UpdateFlow)
	// 	borrows.DELETE("/:accountFlowId/flows/:flowId", borrowController.DeleteFlow)
	// }

	// statistics := g.Group("/statistics")
	// {
	// 	statisticController := statistic.NewStatisticController()
	//
	// 	statistics.GET("", statisticController.List)
	// 	statistics.GET("/years/:year", statisticController.YearStatistic)
	//
	// }
}
