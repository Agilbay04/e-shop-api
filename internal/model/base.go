package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	CreatedBy *uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by"`
	UpdatedAt time.Time	`gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;column:updated_by" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}