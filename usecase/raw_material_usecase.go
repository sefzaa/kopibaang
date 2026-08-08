package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/domain/entity"
)

type RawMaterialUsecase struct {
	materialRepo domain.RawMaterialRepository
}

func NewRawMaterialUsecase(materialRepo domain.RawMaterialRepository) *RawMaterialUsecase {
	return &RawMaterialUsecase{materialRepo}
}

func (u *RawMaterialUsecase) AddMaterial(ctx context.Context, req dto.RawMaterialRequest) error {
	material := &entity.RawMaterial{
		ID:        uuid.New(),
		Name:      req.Name,
		Quantity:  req.Quantity,
		Unit:      req.Unit,
		Price:     req.Price,
		Source:    req.Source,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.materialRepo.Create(ctx, material)
}

func (u *RawMaterialUsecase) UpdateMaterial(ctx context.Context, id string, req dto.RawMaterialRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	existing, err := u.materialRepo.GetByID(ctx, uid)
	if err != nil {
		return errors.New("raw material not found")
	}

	existing.Name = req.Name
	existing.Quantity = req.Quantity
	existing.Unit = req.Unit
	existing.Price = req.Price
	existing.Source = req.Source
	existing.UpdatedAt = time.Now()

	return u.materialRepo.Update(ctx, existing)
}

func (u *RawMaterialUsecase) DeleteMaterial(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	return u.materialRepo.Delete(ctx, uid)
}

func (u *RawMaterialUsecase) GetAllMaterials(ctx context.Context) ([]dto.RawMaterialResponse, error) {
	materials, err := u.materialRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.RawMaterialResponse
	for _, m := range materials {
		responses = append(responses, dto.RawMaterialResponse{
			ID:        m.ID.String(),
			Name:      m.Name,
			Quantity:  m.Quantity,
			Unit:      m.Unit,
			Price:     m.Price,
			Source:    m.Source,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		})
	}

	// Mengembalikan array kosong jika belum ada data agar konsisten
	if responses == nil {
		responses = []dto.RawMaterialResponse{}
	}

	return responses, nil
}