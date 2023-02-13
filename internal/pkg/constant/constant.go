package constant

const (
	// XUserIdKey 在gin.context中表示当前登录用户id的key
	XUserIdKey = "X-User-Id"
	// XRequestIDKey 在gin.context中表示请求id的key
	XRequestIDKey = "X-Request-Id"
	// SqlDBKey 在gin.context中表示数据库连接的key
	SqlDBKey = "sql-db"

	// AccountIdNo 不选择账户时的默认账户id
	AccountIdNo = -1

	// DelFalse 数据库未删除的值
	DelFalse = 0
	// DelTrue 数据库记录删除标志
	DelTrue = 1
)

const (
	// AccountTypeExpense 账户类型 支出
	AccountTypeExpense = iota

	// AccountTypeIncome 账户类型 收入
	AccountTypeIncome

	// AccountTypeTransferIn 账户类型 转入
	AccountTypeTransferIn

	// AccountTypeTransferOut 账户类型 转出
	AccountTypeTransferOut

	// AccountTypeBorrow 账户类型 借入
	AccountTypeBorrow

	// AccountTypeLend 账户类型 借出
	AccountTypeLend

	// AccountTypeModify 账户类型 修改余额
	AccountTypeModify

	// AccountTypeStill 账户类型 还款
	AccountTypeStill

	// AccountTypeHarvest 账户类型 收款
	AccountTypeHarvest
)

// 0-新增，1-修改，2-已同步
const (
	SYNC_ADD = iota
	SYNC_UPDATE
	SYNC_SUCCESS
)
