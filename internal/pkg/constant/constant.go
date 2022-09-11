package constant

const (
	// XUserIdKey 在gin.context中表示当前登录用户id的key
	XUserIdKey = "X-User-Id"
	// XRequestIDKey 在gin.context中表示请求id的key
	XRequestIDKey = "X-Request-Id"
	// SqlDBKey 在gin.context中表示数据库连接的key
	SqlDBKey = "sql-db"

	// BillTypeExpense 账单类型，支出
	BillTypeExpense = 0
	// BillTypeIncome 账单类型，收入
	BillTypeIncome = 1

	// AccountIdNo 不选择账户时的默认账户id
	AccountIdNo = -1
)
