package model

import (
	"gorm.io/gorm"
)

type Payway struct {
	gorm.Model
	Name   string `gorm:"column:name" type:varchar(50);not null;"`
	Config string `gorm:"column:config;type:varchar(255);not null"`
}

func (Payway) TableName() string {
	return "payway"
}

