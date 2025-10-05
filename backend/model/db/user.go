package db

import "time"

// User 用户表
type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(250);not null;uniqueIndex:account_username_UNIQUE;comment:用户名" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(50);not null" json:"password"`
	Salt      string    `gorm:"column:salt;type:varchar(50);not null;default:0;comment:密码盐值" json:"salt"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
