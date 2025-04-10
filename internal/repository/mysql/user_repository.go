package mysql 

import (
	"gorm.io/gorm"
	"payment/internal/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserReposity(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(User *model.User) error {
	return r.DB.Create(User).Error
}

func (r *UserRepository) Update(User *model.User) error {
	return r.DB.Save(User).Error
}