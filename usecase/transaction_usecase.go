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

// ADMIN: Get Order History
func (u *TransactionUsecase) GetAdminOrderHistory(ctx context.Context, req dto.OrderHistoryQueryRequest) (dto.OrderHistoryResponse, error) {
	if req.Page <= 0 { req.Page = 1 }
	if req.Limit <= 0 { req.Limit = 10 }

	// Konfigurasi Zona Waktu St. Petersburg
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", 3*3600) // Fallback ke UTC+3 jika tzdata tidak ada di OS server
	}

	now := time.Now().In(loc)
	var start, end time.Time

	// Logika filter tanggal
	switch req.Filter {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		start = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
		end = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, loc)
	case "this_week":
		weekday := int(now.Weekday())
		if weekday == 0 { weekday = 7 } // Sesuaikan agar Minggu jadi hari ke-7
		startOfWeek := now.AddDate(0, 0, -weekday+1)
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, loc)
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	case "this_month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		// Trick bulan depan tanggal 0 untuk mendapat hari terakhir bulan ini
		end = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, loc)
	case "custom":
		if req.StartDate == "" || req.EndDate == "" {
			return dto.OrderHistoryResponse{}, errors.New("start_date and end_date are required for custom filter")
		}
		parsedStart, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
		if err != nil { return dto.OrderHistoryResponse{}, errors.New("invalid start_date format, use YYYY-MM-DD") }
		parsedEnd, err := time.ParseInLocation("2006-01-02", req.EndDate, loc)
		if err != nil { return dto.OrderHistoryResponse{}, errors.New("invalid end_date format, use YYYY-MM-DD") }
		
		start = parsedStart
		end = time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 999999999, loc)
	}

	orders, total, err := u.txRepo.GetAllOrderHistory(ctx, start, end, req.Page, req.Limit)
	if err != nil {
		return dto.OrderHistoryResponse{}, err
	}

	return u.mapToOrderHistoryResponse(orders, total, req.Page, req.Limit), nil
}

// USER: Get Own Order History
func (u *TransactionUsecase) GetUserOrderHistory(ctx context.Context, userIDStr string, page, limit int) (dto.OrderHistoryResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.OrderHistoryResponse{}, errors.New("invalid user id")
	}

	if page <= 0 { page = 1 }
	if limit <= 0 { limit = 10 }

	orders, total, err := u.txRepo.GetOrderHistoryByUser(ctx, userID, page, limit)
	if err != nil {
		return dto.OrderHistoryResponse{}, err
	}

	return u.mapToOrderHistoryResponse(orders, total, page, limit), nil
}

// Helper untuk Mapping Entity ke DTO
func (u *TransactionUsecase) mapToOrderHistoryResponse(orders []entity.Order, total int64, page, limit int) dto.OrderHistoryResponse {
	var orderRes []dto.OrderResponse
	for _, o := range orders {
		var items []dto.OrderItemDetailResponse
		for _, i := range o.Items {
			items = append(items, dto.OrderItemDetailResponse{
				ProductID:   i.ProductID.String(),
				Quantity:    i.Quantity,
				PriceAtTime: i.PriceAtTime,
			})
		}
		
		userID := ""
		if o.UserID != nil { userID = o.UserID.String() }
		
		orderRes = append(orderRes, dto.OrderResponse{
			OrderID:     o.ID.String(),
			UserID:      userID,
			TotalAmount: o.TotalAmount,
			Discount:    o.Discount,
			FinalAmount: o.FinalAmount,
			IsRedeem:    o.IsRedeem,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
			Items:       items,
		})
	}
	
	if orderRes == nil {
		orderRes = []dto.OrderResponse{} // Cegah JSON bernilai null, jadi [] kosong
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 { totalPages = 1 }

	return dto.OrderHistoryResponse{
		Orders: orderRes,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: total,
		},
	}
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

	for _, item := range req.Items {
		prodID, _ := uuid.Parse(item.ProductID)
		product, err := u.productRepo.GetByID(ctx, prodID)
		if err != nil || !product.IsActive {
			return dto.CreateOrderResponse{}, fmt.Errorf("product %s not found or inactive", item.ProductID)
		}

		itemPrice := product.Price - product.Discount
		if itemPrice < 0 { itemPrice = 0 }
		
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

	totalDiscountOrder := 0
	
	if req.OrderVoucherCode != "" { 
		orderVoucher, err := u.voucherRepo.GetByCode(ctx, req.OrderVoucherCode)
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

	if req.IsRedeem {
		_ = u.txRepo.UpdateUserPoints(ctx, customerID, 100, false)
		_ = u.txRepo.RecordPointTransaction(ctx, &entity.PointTransaction{
			ID: uuid.New(), UserID: customerID, OrderID: &orderID, Type: "redeemed", Points: 100, Description: "Redeem Free Coffee", CreatedAt: time.Now(),
		})
		_ = u.redisRepo.DeleteQRToken(ctx, "redeem_qr", req.RedeemToken)
	}

	receiptText := receiptutils.GenerateReceiptText(order, productNames)

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