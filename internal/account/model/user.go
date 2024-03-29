package model

type User struct {
	Id           int64  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`           // 用户名
	Phone        string `db:"phone" json:"phone"`                 // 手机号
	WxId         string `db:"wx_id" json:"wx_id"`                 // 微信id
	Gender       int    `db:"gender" json:"gender"`               // 性别，0-男，1-女
	Password     string `db:"password" json:"-"`                  // 密码
	Avatar       string `db:"avatar" json:"avatar"`               // 头像
	RegisterTime int64  `db:"register_time" json:"register_time"` // 注册时间
	UpdateTime   int64  `db:"update_time" json:"update_time"`
}
