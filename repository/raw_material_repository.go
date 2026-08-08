package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"kopibang/domain"
	"kopibang/domain/entity"
)

type rawMaterialRepository struct {
	db *gorm.DB
}

func NewRawMaterialRepository(db *gorm.DB) domain.RawMaterialRepository {
	return &rawMaterialRepository{db}
}

func (r *rawMaterialRepository) Create(ctx context.Context, material *entity.RawMaterial) error {
	return r.db.WithContext(ctx).Create(material).Error
}

func (r *rawMaterialRepository) Update(ctx context.Context, material *entity.RawMaterial) error {
	return r.db.WithContext(ctx).Save(material).Error
}

func (r *rawMaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.RawMaterial{}, id).Error
}

func (r *rawMaterialRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RawMaterial, error) {
	var material entity.RawMaterial
	err := r.db.WithContext(ctx).First(&material, id).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *rawMaterialRepository) GetAll(ctx context.Context) ([]entity.RawMaterial, error) {
	var materials []entity.RawMaterial
	err := r.db.WithContext(ctx).Order("updated_at DESC").Find(&materials).Error
	return materials, err
}