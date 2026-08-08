package bootstrap

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func NewFirebaseMessagingClient(env *Env) *messaging.Client {
	ctx := context.Background()

	// Kita paksa secara hardcode agar Firebase langsung membaca file JSON di root folder
	opt := option.WithCredentialsFile("firebase-service-account.json")

	// Tambahkan Config berisi Project ID (ambil dari isi file JSON kamu)
	config := &firebase.Config{
		ProjectID: "kopibang-fdec6",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase App: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase Messaging Client: %v", err)
	}

	log.Println("Firebase Cloud Messaging (FCM) terhubung dengan sukses!")
	return client
}