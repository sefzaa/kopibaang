package entity

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Username     string    `gorm:"type:varchar(255);unique;not null"` // Tambahan baru
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:enum('barista', 'customer');default:'customer'"`
	Points       int       `gorm:"default:0"`
	FCMToken     string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}