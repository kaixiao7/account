package budget

import "kaixiao7/account/internal/account/service"

type BudgetController struct {
	budgetSrv service.BudgetSrv
}

func NewBudgetController() *BudgetController {
	return &BudgetController{budgetSrv: service.NewBudgetSrv()}
}