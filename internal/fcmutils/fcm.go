package fcmutils

import (
	"context"
	"log"

	"firebase.google.com/go/v4/messaging"
)

// SendToTopic digunakan untuk broadcast (Barista Available, Menu Baru, Voucher Baru, Custom Push)
func SendToTopic(client *messaging.Client, topic string, title string, body string) {
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Dijalankan di background agar tidak memblokir response API utama
	go func() {
		response, err := client.Send(context.Background(), message)
		if err != nil {
			log.Printf("Gagal mengirim notifikasi ke topic %s: %v", topic, err)
			return
		}
		log.Printf("Notifikasi berhasil dikirim ke topic %s: %s", topic, response)
	}()
}

// SendToToken digunakan untuk mengirim ke 1 user spesifik (Pengingat Poin)
func SendToToken(client *messaging.Client, token string, title string, body string) {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	go func() {
		_, err := client.Send(context.Background(), message)
		if err != nil {
			log.Printf("Gagal mengirim notifikasi ke token %s: %v", token, err)
		}
	}()
}