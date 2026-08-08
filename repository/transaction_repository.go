package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kopibang/domain"
	"kopibang/domain/entity"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	// Gunakan transaction agar Order dan OrderItem tersimpan berbarengan secara aman
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *transactionRepository) GetOrderHistoryByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *transactionRepository) GetRedeemHistory(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).Preload("Items").Where("is_redeem = ?", true).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *transactionRepository) RecordPointTransaction(ctx context.Context, pt *entity.PointTransaction) error {
	return r.db.WithContext(ctx).Create(pt).Error
}

func (r *transactionRepository) GetPointHistoryByUser(ctx context.Context, userID uuid.UUID) ([]entity.PointTransaction, error) {
	var points []entity.PointTransaction
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&points).Error
	return points, err
}

func (r *transactionRepository) UpdateUserPoints(ctx context.Context, userID uuid.UUID, points int, isAddition bool) error {
	op := "+"
	if !isAddition {
		op = "-"
	}
	// Menggunakan raw expression GORM untuk mencegah race condition saat update poin
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points "+op+" ?", points)).Error
}