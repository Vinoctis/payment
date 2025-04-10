package mysql 

import (
	"gorm.io/gorm"
	"payment/internal/model"
	"errors"
)

type OrderRepository interface {
	Create(order *model.Order) error
	Update(order *model.Order) error
	FindByOrderNo(orderNo string) (*model.Order, error)
	GetPayway(order *model.Order) (model.Payway, error)
}

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.DB.Create(order).Error
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.DB.Save(order).Error
}

func (r *orderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.DB.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetPayway(order *model.Order) (model.Payway, error) {
	if order == nil || order.ID == 0 {
		return model.Payway{}, errors.New("无效的订单")
	}

	err := r.DB.Model(order).Preload("Payway").First(order).Error
	if err!= nil {
		return model.Payway{}, errors.New("获取支付方式失败")
	}

	if order.Payway.ID == 0 {
		return model.Payway{}, errors.New("支付方式不存在")
	}

	return order.Payway, nil
}