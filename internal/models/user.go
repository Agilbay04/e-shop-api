package models

import (
	"e-shop-api/internal/constants"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Username string             `gorm:"unique;not null;column:username" json:"username"`
	Email    string             `gorm:"unique;not null;column:email" json:"email"`
	Password string             `gorm:"not null;column:password" json:"password"`
	Role     constants.UserRole `gorm:"type:varchar(20);default:buyer;column:role" json:"role"`
	IsActive bool               `gorm:"default:true;column:is_active" json:"is_active"`
	Picture  string             `gorm:"type:varchar(500);column:picture" json:"picture"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") && u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}
