package entity

import (
	"time"
	"github.com/google/uuid"
)

type RawMaterial struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Quantity  float64   `gorm:"not null"`
	Unit      string    `gorm:"type:enum('gr', 'kg', 'ml', 'liter', 'pcs', 'pack', 'bottle', 'cup');not null"`
	Price     int       `gorm:"not null"`
	Source    string    `gorm:"type:varchar(255)"` 
	CreatedAt time.Time `gorm:"autoCreateTime"` // Ditambahkan autoCreateTime
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Ditambahkan autoUpdateTime
}