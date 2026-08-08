package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID    `gorm:"type:char(36);primary_key"`
	Name        string       `gorm:"type:varchar(255);not null"`
	Description string       `gorm:"type:text"`
	Price       int          `gorm:"not null"`
	Discount    int          `gorm:"default:0"` // Potongan manual
	VoucherID   *uuid.UUID   `gorm:"type:char(36);index"`
	Voucher     *Voucher     `gorm:"foreignKey:VoucherID"`
	Volume      string       `gorm:"type:varchar(50);not null"`
	ImageURLs   []string     `gorm:"type:json;serializer:json"` // UPDATE: Dukungan banyak foto (MySQL JSON)
	IsActive    bool         `gorm:"default:true"`
	Ingredients []Ingredient `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


type Ingredient struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	ProductID uuid.UUID `gorm:"type:char(36);not null;index"`
	Name      string    `gorm:"type:varchar(255);not null"` // Contoh: "Espresso", "Susu Oat"
	Grammage  string    `gorm:"type:varchar(50);not null"`  // Contoh: "30ml", "15gr"
}