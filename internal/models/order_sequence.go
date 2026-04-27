package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderSequence struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	Date          string    `gorm:"type:date;uniqueIndex;not null;column:date" json:"date"`
	LastSequence  int       `gorm:"type:int;not null;column:last_sequence" json:"last_sequence"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (OrderSequence) TableName() string {
	return "order_sequences"
}

func (os *OrderSequence) BeforeCreate(tx *gorm.DB) (err error) {
	if os.ID == uuid.Nil {
		os.ID = uuid.New()
	}
	return
}