package  mysql 

import (
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Order() OrderRepository
	Payway() PaywayRepository
}

type paymentRepository struct {
	orderRepo OrderRepository
	paywayRepo PaywayRepository
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		orderRepo: NewOrderRepository(db),
		paywayRepo: NewPaywayRepository(db),
	}
}

func (r *paymentRepository) Order() OrderRepository {
	return r.orderRepo
}

func (r *paymentRepository) Payway() PaywayRepository {
	return r.paywayRepo
}