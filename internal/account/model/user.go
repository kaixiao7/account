package model

type User struct {
	Id           int    `gorm:"column:id;primary_key" json:"id"`
	Username     string `gorm:"column:username" json:"username"`           // 用户名
	Phone        string `gorm:"column:phone" json:"phone"`                 // 手机号
	WxId         string `gorm:"column:wx_id" json:"wx_id"`                 // 微信id
	Gender       int    `gorm:"column:gender" json:"gender"`               // 性别，0-男，1-女
	Password     string `gorm:"column:password" json:"-"`                  // 密码
	Avatar       string `gorm:"column:avatar" json:"avatar"`               // 头像
	RegisterTime string `gorm:"column:register_time" json:"register_time"` // 注册时间
	UpdateTime   string `gorm:"column:update_time" json:"update_time"`
}
