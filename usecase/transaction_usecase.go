package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/domain/entity"
	"kopibang/internal/receiptutils"
)

type TransactionUsecase struct {
	txRepo      domain.TransactionRepository
	productRepo domain.ProductRepository
	voucherRepo domain.VoucherRepository
	redisRepo   domain.RedisRepository
}

func NewTransactionUsecase(txRepo domain.TransactionRepository, productRepo domain.ProductRepository, voucherRepo domain.VoucherRepository, redisRepo domain.RedisRepository) *TransactionUsecase {
	return &TransactionUsecase{txRepo, productRepo, voucherRepo, redisRepo}
}

// ADMIN: Input Order
func (u *TransactionUsecase) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error) {
	orderID := uuid.New()
	subtotalAmount := 0
	var orderItems []entity.OrderItem
	productNames := make(map[string]string)

	var customerID uuid.UUID
	if req.IsRedeem {
		if req.RedeemToken == "" { return dto.CreateOrderResponse{}, errors.New("redeem token is required") }
		userIDStr, err := u.redisRepo.GetQRTokenData(ctx, "redeem_qr", req.RedeemToken)
		if err != nil { return dto.CreateOrderResponse{}, errors.New("invalid or expired redeem token") }
		customerID, _ = uuid.Parse(userIDStr)
	}

	// 1. ITERASI ITEM (Diskon Menu + Voucher Menu bertipe 'menu_promo')
	for _, item := range req.Items {
		prodID, _ := uuid.Parse(item.ProductID)
		product, err := u.productRepo.GetByID(ctx, prodID)
		if err != nil || !product.IsActive {
			return dto.CreateOrderResponse{}, fmt.Errorf("product %s not found or inactive", item.ProductID)
		}

		itemPrice := product.Price - product.Discount
		if itemPrice < 0 { itemPrice = 0 }
		
		// PENTING: Hanya potong harga jika voucher aktif dan bertipe 'menu_promo'
		if product.VoucherID != nil && product.Voucher != nil && product.Voucher.IsActive && product.Voucher.Type == "menu_promo" {
			if product.Voucher.DiscountType == "percentage" {
				itemPrice -= (itemPrice * product.Voucher.DiscountValue) / 100
			} else {
				itemPrice -= product.Voucher.DiscountValue
			}
		}

		if itemPrice < 0 { itemPrice = 0 }

		subtotalAmount += itemPrice * item.Quantity
		orderItems = append(orderItems, entity.OrderItem{
			ID:          uuid.New(),
			OrderID:     orderID,
			ProductID:   prodID,
			Quantity:    item.Quantity,
			PriceAtTime: itemPrice,
		})
		productNames[prodID.String()] = product.Name
	}

	// 2. DISKON LEVEL KERANJANG (Bertumpuk dengan menu diskon, khusus 'cart_discount')
	totalDiscountOrder := 0
	
	if req.OrderVoucherCode != "" { 
		orderVoucher, err := u.voucherRepo.GetByCode(ctx, req.OrderVoucherCode)
		// PENTING: Hanya potong total jika valid dan bertipe 'cart_discount'
		if err == nil && orderVoucher.IsActive && subtotalAmount >= orderVoucher.MinPurchase && orderVoucher.Type == "cart_discount" {
			if orderVoucher.DiscountType == "percentage" {
				totalDiscountOrder = (subtotalAmount * orderVoucher.DiscountValue) / 100
			} else {
				totalDiscountOrder = orderVoucher.DiscountValue
			}
		}
	}
	
	finalAmount := subtotalAmount - totalDiscountOrder
	if finalAmount < 0 { finalAmount = 0 }

	order := &entity.Order{
		ID:          orderID,
		TotalAmount: subtotalAmount,
		Discount:    totalDiscountOrder,
		FinalAmount: finalAmount,
		IsRedeem:    req.IsRedeem,
		Status:      "completed",
		Items:       orderItems,
		CreatedAt:   time.Now(),
	}

	if req.IsRedeem { order.UserID = &customerID }

	if err := u.txRepo.CreateOrder(ctx, order); err != nil {
		return dto.CreateOrderResponse{}, err
	}

	// 3. PROSES POIN
	if req.IsRedeem {
		_ = u.txRepo.UpdateUserPoints(ctx, customerID, 100, false)
		_ = u.txRepo.RecordPointTransaction(ctx, &entity.PointTransaction{
			ID: uuid.New(), UserID: customerID, OrderID: &orderID, Type: "redeemed", Points: 100, Description: "Redeem Free Coffee", CreatedAt: time.Now(),
		})
		_ = u.redisRepo.DeleteQRToken(ctx, "redeem_qr", req.RedeemToken)
	}

	// Buat teks struk WhatsApp
	receiptText := receiptutils.GenerateReceiptText(order, productNames)

	// 4. GENERATE EARN QR
	var earnToken string
	if !req.IsRedeem {
		earnToken = uuid.New().String()
		redisData := fmt.Sprintf("%s|%d", orderID.String(), finalAmount)
		_ = u.redisRepo.SaveQRToken(ctx, "earn_qr", earnToken, redisData, 10*time.Minute)
	}

	return dto.CreateOrderResponse{
		OrderID:        orderID.String(),
		FinalAmount:    finalAmount,
		EarnPointToken: earnToken,
		ReceiptText:    receiptText,
	}, nil
}

func (u *TransactionUsecase) GenerateRedeemQR(ctx context.Context, userIDStr string) (dto.RedeemQRResponse, error) {
	token := uuid.New().String()
	_ = u.redisRepo.SaveQRToken(ctx, "redeem_qr", token, userIDStr, 5*time.Minute)
	return dto.RedeemQRResponse{RedeemToken: token}, nil
}

func (u *TransactionUsecase) ScanEarnQR(ctx context.Context, userIDStr string, req dto.ScanEarnPointRequest) error {
	userID, _ := uuid.Parse(userIDStr)

	_, err := u.redisRepo.GetQRTokenData(ctx, "earn_qr", req.EarnToken)
	if err != nil { return errors.New("QR code expired or invalid") }

	pointsEarned := 20 
	
	err = u.txRepo.UpdateUserPoints(ctx, userID, pointsEarned, true)
	if err != nil { return err }

	_ = u.txRepo.RecordPointTransaction(ctx, &entity.PointTransaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        "earned",
		Points:      pointsEarned,
		Description: "Points from purchase",
		CreatedAt:   time.Now(),
	})

	_ = u.redisRepo.DeleteQRToken(ctx, "earn_qr", req.EarnToken)
	return nil
}