package repository

import (
	"e-shop-api/internal/model"
	"gorm.io/gorm"
)

type UserQueryRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) UserQueryRepository {
	return &userQueryRepository{db}
}

func (r *userQueryRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userQueryRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}