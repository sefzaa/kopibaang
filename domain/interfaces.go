package domain

import (
	"context"
	"github.com/google/uuid"
	"kopibang/domain/entity" // Sesuaikan nama modul
	"kopibang/domain/dto"
	"time"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error
	Update(ctx context.Context, user *entity.User) error // Tambahan baru
}

type RedisRepository interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiresIn time.Duration) error
	ValidateRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string) error
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string) error

	SaveOTP(ctx context.Context, email string, otp string, expiresIn time.Duration) error
	ValidateOTP(ctx context.Context, email string, otp string) error
	DeleteOTP(ctx context.Context, email string) error

	SaveQRToken(ctx context.Context, tokenType string, tokenID string, data string, expiresIn time.Duration) error
	GetQRTokenData(ctx context.Context, tokenType string, tokenID string) (string, error)
	DeleteQRToken(ctx context.Context, tokenType string, tokenID string) error

	SetState(ctx context.Context, key string, value string) error
	GetState(ctx context.Context, key string) (string, error)

	SaveResetToken(ctx context.Context, email string, token string, expiresIn time.Duration) error
	ValidateResetToken(ctx context.Context, email string, token string) error
	DeleteResetToken(ctx context.Context, email string) error
}

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetAll(ctx context.Context, onlyActive bool) ([]entity.Product, error)
}

type TransactionRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderHistoryByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	GetRedeemHistory(ctx context.Context) ([]entity.Order, error)
	
	RecordPointTransaction(ctx context.Context, pt *entity.PointTransaction) error
	GetPointHistoryByUser(ctx context.Context, userID uuid.UUID) ([]entity.PointTransaction, error)
	
	UpdateUserPoints(ctx context.Context, userID uuid.UUID, points int, isAddition bool) error
}

type DashboardRepository interface {
	GetTotalSales(ctx context.Context, start time.Time, end time.Time) (int, error)
	GetTotalOrders(ctx context.Context, start time.Time, end time.Time) (int, error)
	GetTotalPointsRedeemed(ctx context.Context, start time.Time, end time.Time) (int, error)
	GetMenuStats(ctx context.Context, start time.Time, end time.Time) ([]dto.MenuStatResponse, error)
}

type RawMaterialRepository interface {
	Create(ctx context.Context, material *entity.RawMaterial) error
	Update(ctx context.Context, material *entity.RawMaterial) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RawMaterial, error)
	GetAll(ctx context.Context) ([]entity.RawMaterial, error)
}

type SettingRepository interface {
	GetSetting(ctx context.Context, key string) (string, error)
	UpdateSetting(ctx context.Context, key string, value string) error
}



type VoucherRepository interface {
	Create(ctx context.Context, voucher *entity.Voucher) error
	Update(ctx context.Context, voucher *entity.Voucher) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Voucher, error)
	GetByCode(ctx context.Context, code string) (*entity.Voucher, error)
	GetAll(ctx context.Context) ([]entity.Voucher, error)
}