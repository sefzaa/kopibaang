package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kopibang/domain"
	"kopibang/domain/entity"
)

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) domain.VoucherRepository {
	return &voucherRepository{db}
}

func (r *voucherRepository) Create(ctx context.Context, voucher *entity.Voucher) error {
	return r.db.WithContext(ctx).Create(voucher).Error
}

func (r *voucherRepository) Update(ctx context.Context, voucher *entity.Voucher) error {
	return r.db.WithContext(ctx).Save(voucher).Error
}

func (r *voucherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Voucher{}, id).Error
}

func (r *voucherRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.WithContext(ctx).First(&voucher, id).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) GetByCode(ctx context.Context, code string) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) GetAll(ctx context.Context) ([]entity.Voucher, error) {
	var vouchers []entity.Voucher
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&vouchers).Error
	return vouchers, err
}