package entity

import (
	"time"
	"github.com/google/uuid"
)

type Voucher struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	Code          string    `gorm:"type:varchar(50);unique;not null"`
	Type          string    `gorm:"type:enum('menu_promo', 'cart_discount');default:'cart_discount'"` // TAMBAHAN: Untuk membedakan jenis voucher
	DiscountType  string    `gorm:"type:enum('nominal', 'percentage');default:'nominal'"`
	DiscountValue int       `gorm:"not null"` 
	MinPurchase   int       `gorm:"default:0"`
	IsActive      bool      `gorm:"default:true"`
	StartDate     time.Time `gorm:"not null"`
	EndDate       time.Time `gorm:"not null"`
}