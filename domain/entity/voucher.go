package entity

import (
	"time"
	"github.com/google/uuid"
)

type Voucher struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key"`
	Code           string    `gorm:"type:varchar(50);unique;not null"`
	Type           string    `gorm:"type:enum('menu_promo', 'cart_discount');default:'cart_discount'"` 
	
	// Sesuaikan nama variabel dengan kolom di database
	DiscountAmount int       `gorm:"column:discount_amount;not null"` 
	MinPurchase    int       `gorm:"column:min_purchase;default:0"`
	IsActive       bool      `gorm:"column:is_active;default:true"`
	
	StartDate      time.Time `gorm:"column:start_date;not null"`
	EndDate        time.Time `gorm:"column:end_date;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}