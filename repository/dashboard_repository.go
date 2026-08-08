package repository

import (
	"context"
	"time"
	"gorm.io/gorm"
	"kopibang/domain"
	"kopibang/domain/dto"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) GetTotalSales(ctx context.Context, start time.Time, end time.Time) (int, error) {
	var totalSales *int
	err := r.db.WithContext(ctx).Table("orders").
		Select("SUM(final_amount)").
		Where("status = ?", "completed").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Scan(&totalSales).Error

	if err != nil || totalSales == nil {
		return 0, err
	}
	return *totalSales, nil
}

func (r *dashboardRepository) GetTotalOrders(ctx context.Context, start time.Time, end time.Time) (int, error) {
	var totalOrders int64
	err := r.db.WithContext(ctx).Table("orders").
		Where("status = ?", "completed").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Count(&totalOrders).Error

	return int(totalOrders), err
}

func (r *dashboardRepository) GetTotalPointsRedeemed(ctx context.Context, start time.Time, end time.Time) (int, error) {
	var totalRedeemed int64
	err := r.db.WithContext(ctx).Table("orders").
		Where("is_redeem = ? AND status = ?", true, "completed").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Count(&totalRedeemed).Error

	return int(totalRedeemed), err
}

func (r *dashboardRepository) GetMenuStats(ctx context.Context, start time.Time, end time.Time) ([]dto.MenuStatResponse, error) {
	var menuStats []dto.MenuStatResponse
	
	// Menampilkan SEMUA produk yang terjual pada rentang waktu tersebut (tanpa Limit)
	err := r.db.WithContext(ctx).Table("order_items").
		Select("products.name as product_name, SUM(order_items.quantity) as total_sold").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status = ?", "completed").
		Where("orders.created_at >= ? AND orders.created_at <= ?", start, end).
		Group("products.id, products.name").
		Order("total_sold DESC").
		Scan(&menuStats).Error

	return menuStats, err
}