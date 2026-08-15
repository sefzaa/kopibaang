package entity

import (
	"time"
	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:char(36);primary_key"`
	UserID      *uuid.UUID  `gorm:"type:char(36);index"` // Null jika bukan user aplikasi
	TotalAmount int         `gorm:"not null"`
	Discount    int         `gorm:"default:0"`
	FinalAmount int         `gorm:"not null"`
	IsRedeem    bool        `gorm:"default:false"`
	Status      string      `gorm:"type:enum('pending', 'completed', 'cancelled');default:'completed'"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time
}


type OrderItem struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	OrderID     uuid.UUID `gorm:"type:char(36);not null"`
	ProductID   uuid.UUID `gorm:"type:char(36);not null"`
	Quantity    int       `gorm:"not null"`
	PriceAtTime int       `gorm:"not null"`
	Product     Product   `gorm:"foreignKey:ProductID"` // <--- TAMBAHKAN BARIS INI
}

type PointTransaction struct {
	ID          uuid.UUID  `gorm:"type:char(36);primary_key"`
	UserID      uuid.UUID  `gorm:"type:char(36);not null"`
	OrderID     *uuid.UUID `gorm:"type:char(36)"`
	Type        string     `gorm:"type:enum('earned', 'redeemed');not null"`
	Points      int        `gorm:"not null"`
	Description string     `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
}