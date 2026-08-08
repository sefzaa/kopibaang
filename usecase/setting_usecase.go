package usecase

import (
	"context"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/internal/fcmutils"
	"firebase.google.com/go/v4/messaging"
)

type SettingUsecase struct {
	settingRepo domain.SettingRepository
	redisRepo   domain.RedisRepository
	fcmClient   *messaging.Client // FCM Client ditambahkan agar tidak error
}

func NewSettingUsecase(settingRepo domain.SettingRepository, redisRepo domain.RedisRepository, fcmClient *messaging.Client) *SettingUsecase {
	return &SettingUsecase{settingRepo, redisRepo, fcmClient}
}

func (u *SettingUsecase) UpdateBaristaStatus(ctx context.Context, req dto.BaristaStatusRequest) error {
	statusStr := "offline"
	if *req.IsAvailable {
		statusStr = "online"
	}

	// 1. Simpan ke MySQL sebagai backup permanen
	err := u.settingRepo.UpdateSetting(ctx, "barista_status", statusStr)
	if err != nil {
		return err
	}

	// 2. Simpan ke Redis agar akses pelanggan sangat cepat
	_ = u.redisRepo.SetState(ctx, "system:barista_status", statusStr)

	// Note: Jika kamu nanti implementasi WebSockets, ini adalah tempat yang tepat 
	// untuk melakukan trigger / broadcast ke semua user yang sedang membuka aplikasi.

	if *req.IsAvailable {
        // Trigger notifikasi jika status diubah menjadi online
        fcmutils.SendToTopic(u.fcmClient, "all_users", "kopibang Coffee Buka! ☕", "Barista lagi di kamar nih, yuk pesan kopi sekarang sebelum tutup lagi!")
    }

	return nil
}

func (u *SettingUsecase) GetBaristaStatus(ctx context.Context) (dto.BaristaStatusResponse, error) {
	// Coba ambil dari Redis dulu (karena jauh lebih cepat)
	statusStr, err := u.redisRepo.GetState(ctx, "system:barista_status")
	
	if err != nil || statusStr == "" {
		// Fallback: Jika Redis kosong/mati, ambil dari MySQL
		statusStr, _ = u.settingRepo.GetSetting(ctx, "barista_status")
	}

	isAvailable := false
	if statusStr == "online" {
		isAvailable = true
	}

	return dto.BaristaStatusResponse{
		IsAvailable: isAvailable,
		StatusText:  statusStr,
	}, nil
}