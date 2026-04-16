package model

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name    string    `gorm:"not null;column:name" json:"name"`
	Slug    string    `gorm:"uniqueIndex;column:slug" json:"slug"`
	Price   int       `gorm:"not null;column:price" json:"price"`
	Stock   int       `gorm:"default:0;column:stock" json:"stock"`
	StoreID uuid.UUID `gorm:"type:uuid;column:store_id" json:"store_id"`
	Store   Store     `gorm:"foreignKey:StoreID" json:"store"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Base.BeforeCreate(tx) 

	// Generate Slug: "Baju Koko Pria" -> "baju-koko-pria"
	p.Slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
	return
}