package repository

import (
	"context"
	"time"

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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// Baru: Method untuk History Admin dengan rentang tanggal dan pagination
func (r *transactionRepository) GetAllOrderHistory(ctx context.Context, start time.Time, end time.Time, page int, limit int) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Order{})

	// Hanya aplikasikan filter jika tanggal tidak kosong (misal: filter 'all')
	if !start.IsZero() && !end.IsZero() {
		query = query.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	// Hitung total data keseluruhan (sebelum pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Order("created_at DESC").Limit(limit).Offset(offset).Find(&orders).Error

	return orders, total, err
}

// Diubah: Ditambahkan Pagination dan menghitung Total
func (r *transactionRepository) GetOrderHistoryByUser(ctx context.Context, userID uuid.UUID, page int, limit int) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Order{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Order("created_at DESC").Limit(limit).Offset(offset).Find(&orders).Error
	
	return orders, total, err
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
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("points", gorm.Expr("points "+op+" ?", points)).Error
}