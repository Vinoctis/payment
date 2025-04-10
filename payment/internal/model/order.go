package model

import (
	"gorm.io/gorm"
)
const (
	OrderStatusPending = 0
	OrderStatusWaitPay = 1
	OrderStatusPaid = 2
	OrderStatusConfirmed = 4
)

type Order struct {
	gorm.Model
	UserID int `gorm:"column:user_id;type:int(11);not null;comment:用户ID"`
	OrderNo string `gorm:"uniqueIndex:idx_order_no,length:64;column:order_no;type:varchar(64) not null;comment:订单号"`
	Amount int `gorm:"column:amount;type:int(11);not null;comment:订单金额"`
	PlatformOrderNo string `gorm:"column:platform_order_no;type:varchar(64) not null;comment:平台订单号"`
	PaywayId int `gorm:"column:payway_id;type:int(11);not null;comment:支付方式"`
	Status int32 `gorm:"column:status;default:1;comment:订单状态;type:tinyint"`
	Payway Payway `gorm:"foreignKey:PaywayId"`
}

func (Order) TableName() string {
	return "order"
}
