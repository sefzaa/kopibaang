package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"kopibang/domain"
	"kopibang/domain/entity"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	// Gunakan Transaction untuk memastikan Ingredient lama terhapus dan terganti baru dengan aman
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hapus bahan-bahan lama
		if err := tx.Where("product_id = ?", product.ID).Delete(&entity.Ingredient{}).Error; err != nil {
			return err
		}

		// Update product (akan otomatis insert ingredient baru dari struct product)
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}



func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	// Tambahkan Preload Voucher
	err := r.db.WithContext(ctx).Preload("Ingredients").Preload("Voucher").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(ctx context.Context, onlyActive bool) ([]entity.Product, error) {
	var products []entity.Product
	// Tambahkan Preload Voucher
	query := r.db.WithContext(ctx).Preload("Ingredients").Preload("Voucher")
	
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Find(&products).Error
	return products, err
}