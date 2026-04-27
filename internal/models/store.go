package models

import "github.com/google/uuid"

type Store struct {
	Base
	Name        string    `gorm:"unique;not null;column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	IsActive    bool      `gorm:"default:true;index;column:is_active" json:"is_active"`
	UserID      uuid.UUID `gorm:"type:uuid;index;column:user_id" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
}

func (Store) TableName() string {
	return "stores"
}