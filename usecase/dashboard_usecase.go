package usecase

import (
	"context"
	"errors"
	"time"
	"kopibang/domain"
	"kopibang/domain/dto"
)

type DashboardUsecase struct {
	dashboardRepo domain.DashboardRepository
}

func NewDashboardUsecase(dashboardRepo domain.DashboardRepository) *DashboardUsecase {
	return &DashboardUsecase{dashboardRepo}
}

func (u *DashboardUsecase) GetDashboardData(ctx context.Context, req dto.DashboardRequest) (dto.DashboardResponse, error) {
	start, end, err := u.parseDateRange(req)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	var response dto.DashboardResponse

	response.TotalSales, err = u.dashboardRepo.GetTotalSales(ctx, start, end)
	if err != nil { return response, err }

	response.TotalOrders, err = u.dashboardRepo.GetTotalOrders(ctx, start, end)
	if err != nil { return response, err }

	response.TotalPointsRedeemed, err = u.dashboardRepo.GetTotalPointsRedeemed(ctx, start, end)
	if err != nil { return response, err }

	response.MenuStats, err = u.dashboardRepo.GetMenuStats(ctx, start, end)
	if err != nil { return response, err }

	if response.MenuStats == nil {
		response.MenuStats = []dto.MenuStatResponse{}
	}

	return response, nil
}

// Helper function untuk menentukan rentang waktu
func (u *DashboardUsecase) parseDateRange(req dto.DashboardRequest) (time.Time, time.Time, error) {
	now := time.Now()
	var start, end time.Time

	switch req.Filter {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		start = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, now.Location())
	case "this_week":
		// Mengasumsikan minggu dimulai hari Senin
		offset := int(time.Monday - now.Weekday())
		if offset > 0 {
			offset = -6
		}
		startOfWeek := now.AddDate(0, 0, offset)
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "this_month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "this_year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "custom":
		if req.StartDate == "" || req.EndDate == "" {
			return start, end, errors.New("start_date and end_date are required for custom filter")
		}
		parsedStart, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return start, end, errors.New("invalid start_date format, use YYYY-MM-DD")
		}
		parsedEnd, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return start, end, errors.New("invalid end_date format, use YYYY-MM-DD")
		}
		start = time.Date(parsedStart.Year(), parsedStart.Month(), parsedStart.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 999999999, now.Location())
	default:
		// Default ke "today" jika filter kosong atau tidak dikenali
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	}

	return start, end, nil
}