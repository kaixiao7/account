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
	// AssetTypeExpense 资产类型 支出
	AssetTypeExpense = iota

	// AssetTypeIncome 资产类型 收入
	AssetTypeIncome

	// AssetTypeTransferIn 资产类型 转入
	AssetTypeTransferIn

	// AssetTypeTransferOut 资产类型 转出
	AssetTypeTransferOut

	// AssetTypeBorrowIn 资产类型 借入
	AssetTypeBorrowIn

	// AssetTypeBorrowOut 资产类型 借出
	AssetTypeBorrowOut

	// AssetTypeModify 资产类型 修改余额
	AssetTypeModify
)
