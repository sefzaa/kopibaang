package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"kopibang/bootstrap"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/domain/entity"
	"kopibang/internal/fcmutils"
)

type ProductUsecase struct {
	productRepo domain.ProductRepository
	fcmClient   *messaging.Client
	minioClient *minio.Client 
	env         *bootstrap.Env 
}

func NewProductUsecase(productRepo domain.ProductRepository, fcmClient *messaging.Client, minioClient *minio.Client, env *bootstrap.Env) *ProductUsecase {
	return &ProductUsecase{productRepo, fcmClient, minioClient, env}
}

func (u *ProductUsecase) UploadImageToMinIO(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	bucketName := u.env.MinioBucketName

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("menu/%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = u.minioClient.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("gagal upload ke MinIO: %v", err)
	}

	publicURL := fmt.Sprintf("http://72.60.198.166:9005/%s/%s", bucketName, objectName)
	
	return publicURL, nil
}

func (u *ProductUsecase) CreateMenu(ctx context.Context, req dto.ProductRequest) error {
	productID := uuid.New()
	
	var voucherUUID *uuid.UUID
	if req.VoucherID != nil && *req.VoucherID != "" {
		parsed, err := uuid.Parse(*req.VoucherID)
		if err != nil { return errors.New("invalid voucher ID format") }
		voucherUUID = &parsed
	}

	ingredients := make([]entity.Ingredient, 0, len(req.Ingredients))
	for _, ing := range req.Ingredients {
		ingredients = append(ingredients, entity.Ingredient{
			ID:        uuid.New(),
			ProductID: productID,
			Name:      ing.Name,
			Grammage:  ing.Grammage,
		})
	}

	product := &entity.Product{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Discount:    req.Discount,
		VoucherID:   voucherUUID,
		Volume:      req.Volume,
		ImageURLs:   req.ImageURLs, 
		IsActive:    *req.IsActive,
		Ingredients: ingredients,
	}

	err := u.productRepo.Create(ctx, product)
	if err != nil {
		return err
	}

	if *req.IsActive && u.fcmClient != nil {
		fcmutils.SendToTopic(u.fcmClient, "all_users", "Menu Baru: "+req.Name+" 🎉", "Cek aplikasi sekarang buat cobain menu baru kita!")
	}

	return nil
}

func (u *ProductUsecase) UpdateMenu(ctx context.Context, productID string, req dto.ProductRequest) error {
	id, err := uuid.Parse(productID)
	if err != nil { return errors.New("invalid product ID format") }

	existing, err := u.productRepo.GetByID(ctx, id)
	if err != nil { return errors.New("product not found") }

	var voucherUUID *uuid.UUID
	if req.VoucherID != nil && *req.VoucherID != "" {
		parsed, err := uuid.Parse(*req.VoucherID)
		if err != nil { return errors.New("invalid voucher ID format") }
		voucherUUID = &parsed
	}

	ingredients := make([]entity.Ingredient, 0, len(req.Ingredients))
	for _, ing := range req.Ingredients {
		ingredients = append(ingredients, entity.Ingredient{
			ID:        uuid.New(),
			ProductID: existing.ID,
			Name:      ing.Name,
			Grammage:  ing.Grammage,
		})
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Price = req.Price
	existing.Discount = req.Discount
	existing.VoucherID = voucherUUID
	existing.Volume = req.Volume
	existing.ImageURLs = req.ImageURLs
	existing.IsActive = *req.IsActive
	existing.Ingredients = ingredients

	return u.productRepo.Update(ctx, existing)
}

func (u *ProductUsecase) ToggleMenuStatus(ctx context.Context, productID string) error {
	id, err := uuid.Parse(productID)
	if err != nil { return errors.New("invalid product ID format") }

	existing, err := u.productRepo.GetByID(ctx, id)
	if err != nil { return errors.New("product not found") }

	existing.IsActive = !existing.IsActive 
	return u.productRepo.Update(ctx, existing)
}

func (u *ProductUsecase) DeleteMenu(ctx context.Context, productID string) error {
	id, err := uuid.Parse(productID)
	if err != nil { return errors.New("invalid product ID format") }
	return u.productRepo.Delete(ctx, id)
}

func (u *ProductUsecase) GetMenus(ctx context.Context, role string) ([]dto.ProductResponse, error) {
	onlyActive := role == "customer"
	
	products, err := u.productRepo.GetAll(ctx, onlyActive)
	if err != nil { return nil, err }

	var responses []dto.ProductResponse
	for _, p := range products {
		responses = append(responses, mapProductToResponse(p))
	}
	
	if responses == nil { responses = []dto.ProductResponse{} }
	return responses, nil
}

func (u *ProductUsecase) GetMenuByID(ctx context.Context, productID string, role string) (*dto.ProductResponse, error) {
	id, err := uuid.Parse(productID)
	if err != nil { return nil, errors.New("invalid product ID format") }

	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil { return nil, errors.New("product not found") }

	if role == "customer" && !product.IsActive {
		return nil, errors.New("product not found or archived")
	}

	resp := mapProductToResponse(*product)
	return &resp, nil
}

func mapProductToResponse(p entity.Product) dto.ProductResponse {
	var ings []dto.IngredientResponse
	for _, i := range p.Ingredients {
		ings = append(ings, dto.IngredientResponse{
			ID:       i.ID.String(),
			Name:     i.Name,
			Grammage: i.Grammage,
		})
	}

	var voucherIDStr *string
	voucherCode := ""
	voucherDiscountAmount := 0

	itemPrice := p.Price - p.Discount
	if itemPrice < 0 { itemPrice = 0 }
	
	// Kalkulasi Voucher Khusus Menu (Diperbarui)
	if p.VoucherID != nil && p.Voucher != nil && p.Voucher.IsActive {
		vid := p.VoucherID.String()
		voucherIDStr = &vid
		voucherCode = p.Voucher.Code
		
		// Langsung ambil potongan harga
		voucherDiscountAmount = p.Voucher.DiscountAmount
	}

	finalPrice := itemPrice - voucherDiscountAmount
	if finalPrice < 0 { finalPrice = 0 }

	return dto.ProductResponse{
		ID:              p.ID.String(),
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price,
		Discount:        p.Discount,
		VoucherID:       voucherIDStr,
		VoucherCode:     voucherCode,
		VoucherDiscount: voucherDiscountAmount,
		FinalPrice:      finalPrice,
		Volume:          p.Volume,
		ImageURLs:       p.ImageURLs,
		IsActive:        p.IsActive,
		Ingredients:     ings,
	}
}