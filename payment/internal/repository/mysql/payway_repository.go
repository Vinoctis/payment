package mysql 

import (
	"gorm.io/gorm"
	"payment/internal/model"
)

type PaywayRepository interface {
	Create(Payway *model.Payway) error
	Update(Payway *model.Payway) error
	FindByName(name string) (*model.Payway, error)
}

type paywayRepository struct {
	DB *gorm.DB
}

func NewPaywayRepository(db *gorm.DB) PaywayRepository {
	return &paywayRepository{
		DB: db,
	}
}

func (r *paywayRepository) Create(Payway *model.Payway) error {
	return r.DB.Create(Payway).Error
}

func (r *paywayRepository) Update(Payway *model.Payway) error {
	return r.DB.Save(Payway).Error
}

func (r *paywayRepository) FindByName(name string) (*model.Payway, error){
	var payway *model.Payway
	err := r.DB.Where("name = ?", name).First(payway).Error
	if err != nil {
		return nil, err
	}
	return payway, nil
}