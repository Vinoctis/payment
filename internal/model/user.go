package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string   `gorm:"column:username;type:varchar(64);not null;comment:用户名"`
	Password string   `gorm:"column:password;type:varchar(64);not null;comment:密码"`
	Salt     string   `gorm:"column:salt;type:varchar(64);not null;comment:盐"`
	Email    string   `gorm:"column:email;type:varchar(64);not null;comment:邮箱"`
	Phone    string   `gorm:"column:phone;type:varchar(64);not null;comment:手机号"`
	Status   int      `gorm:"column:status;type:tinyint;not null;comment:状态"`
}

func (User) TableName() string {
	return "user"
}