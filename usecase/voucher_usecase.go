package usecase

import (
	"context"
	"errors"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/domain/entity"
	"kopibang/internal/fcmutils"
)

type VoucherUsecase struct {
	voucherRepo domain.VoucherRepository
	fcmClient   *messaging.Client
}

func NewVoucherUsecase(voucherRepo domain.VoucherRepository, fcmClient *messaging.Client) *VoucherUsecase {
	return &VoucherUsecase{voucherRepo, fcmClient}
}

func (u *VoucherUsecase) CreateVoucher(ctx context.Context, req dto.VoucherRequest) error {
	_, err := u.voucherRepo.GetByCode(ctx, req.Code)
	if err == nil {
		return errors.New("voucher code already exists")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil { return errors.New("invalid start_date format, use YYYY-MM-DD") }

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil { return errors.New("invalid end_date format, use YYYY-MM-DD") }

	if startDate.After(endDate) {
		return errors.New("start_date cannot be after end_date")
	}

	voucher := &entity.Voucher{
		ID:             uuid.New(),
		Code:           req.Code,
		Type:           req.Type, 
		DiscountAmount: req.DiscountAmount, // PERBAIKAN: Disesuaikan dengan Entity
		MinPurchase:    req.MinPurchase,
		IsActive:       *req.IsActive,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	err = u.voucherRepo.Create(ctx, voucher)
	if err != nil { return err }

	if *req.IsActive && u.fcmClient != nil && req.Type == "cart_discount" {
		fcmutils.SendToTopic(u.fcmClient, "all_users", "Promo Spesial! Gunakan kode: "+req.Code, "Ada potongan harga baru nih. Buruan sikat!")
	}

	return nil
}

func (u *VoucherUsecase) UpdateVoucher(ctx context.Context, id string, req dto.VoucherRequest) error {
	uid, err := uuid.Parse(id)
	if err != nil { return errors.New("invalid voucher ID format") }

	existing, err := u.voucherRepo.GetByID(ctx, uid)
	if err != nil { return errors.New("voucher not found") }

	if existing.Code != req.Code {
		_, errCheck := u.voucherRepo.GetByCode(ctx, req.Code)
		if errCheck == nil { return errors.New("voucher code already in use") }
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil { return errors.New("invalid start_date format") }
	
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil { return errors.New("invalid end_date format") }

	existing.Code = req.Code
	existing.Type = req.Type 
	existing.DiscountAmount = req.DiscountAmount // PERBAIKAN: Disesuaikan dengan Entity
	existing.MinPurchase = req.MinPurchase
	existing.IsActive = *req.IsActive
	existing.StartDate = startDate
	existing.EndDate = endDate

	return u.voucherRepo.Update(ctx, existing)
}

func (u *VoucherUsecase) DeleteVoucher(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return errors.New("invalid voucher ID format") }
	return u.voucherRepo.Delete(ctx, uid)
}

func (u *VoucherUsecase) GetAllVouchers(ctx context.Context) ([]dto.VoucherResponse, error) {
	vouchers, err := u.voucherRepo.GetAll(ctx)
	if err != nil { return nil, err }

	var responses []dto.VoucherResponse
	for _, v := range vouchers {
		responses = append(responses, dto.VoucherResponse{
			ID:             v.ID.String(),
			Code:           v.Code,
			Type:           v.Type, 
			DiscountAmount: v.DiscountAmount, // PERBAIKAN: Disesuaikan dengan Entity
			MinPurchase:    v.MinPurchase,
			IsActive:       v.IsActive,
			StartDate:      v.StartDate.Format("2006-01-02"),
			EndDate:        v.EndDate.Format("2006-01-02"),
		})
	}
	
	if responses == nil { responses = []dto.VoucherResponse{} }

	return responses, nil
}