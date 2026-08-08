package scheduler

import (
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"kopibang/domain/entity"
	"kopibang/internal/fcmutils"
)

func StartPointReminderCron(db *gorm.DB, fcmClient *messaging.Client) {
	// Inisialisasi Cron
	c := cron.New()

	// Jadwal: Setiap hari Jumat jam 16:00 (Sore hari, waktu pas buat ngopi)
	// Format Cron standar: "Menit Jam Tanggal Bulan Hari"
	_, err := c.AddFunc("0 16 * * 5", func() {
		log.Println("Menjalankan cron job pengingat poin...")
		
		var users []entity.User
		// Ambil user yang punya fcm_token, dan poinnya antara 50 sampai 99
		err := db.Where("fcm_token IS NOT NULL AND fcm_token != '' AND points >= 50 AND points < 100").Find(&users).Error
		if err != nil {
			log.Printf("Gagal mengambil data user untuk cron: %v", err)
			return
		}

		for _, user := range users {
			sisaPoin := 100 - user.Points
			title := "Kopi Gratis Hampir di Tangan! ☕"
			body := fmt.Sprintf("Hai %s, kumpulin %d poin lagi buat dapetin gratis kopi. Yuk order sekarang!", user.Name, sisaPoin)
			
			fcmutils.SendToToken(fcmClient, user.FCMToken, title, body)
		}
	})

	if err != nil {
		log.Fatalf("Gagal mendaftarkan cron job: %v", err)
	}

	c.Start()
	log.Println("Cron job penjadwalan notifikasi aktif.")
}