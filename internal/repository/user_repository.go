package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tx *gorm.DB, user *model.User) error
	Update(tx *gorm.DB, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(tx *gorm.DB, user *model.User) error {
	if tx != nil {
		return tx.Create(user).Error
	}
	return r.db.Create(user).Error
}

func (r *userRepository) Update(tx *gorm.DB, user *model.User) error {
	if tx != nil {
		return tx.Save(user).Error
	}
	return r.db.Save(user).Error
}
