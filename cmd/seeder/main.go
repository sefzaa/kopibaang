package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"kopibang/bootstrap" // Sesuaikan dengan nama modulmu
)

func main() {
	// Inisialisasi koneksi DB dari bootstrap yang sudah ada
	app := bootstrap.App()
	defer app.CloseDBConnection()

	email := "admin@sefza.com"
	password := "rahasia123"
	name := "Sefza"
	username := "admin_sefza" // Tambahan untuk skema baru

	// Hashing password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal melakukan hash password:", err)
	}

	// Generate UUID dan waktu saat ini
	id := uuid.New().String()
	now := time.Now()

	// Masukkan ke database dengan role 'barista'
	query := `INSERT INTO users (id, name, username, email, password_hash, role, points, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// TAMBAHKAN .Error DI BELAKANGNYA
	err = app.DB.Exec(query, id, name, username, email, string(hashedPassword), "barista", 0, now, now).Error
	if err != nil {
		log.Fatal("Gagal menjalankan seeder (mungkin email atau username sudah digunakan):", err)
	}

	fmt.Println("==================================================")
	fmt.Println("✅ Seeder berhasil! Akun admin telah dibuat.")
	fmt.Printf("   Email    : %s\n", email)
	fmt.Printf("   Password : %s\n", password)
	fmt.Println("==================================================")
}