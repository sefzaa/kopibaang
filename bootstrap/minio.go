package bootstrap

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(env *Env) *minio.Client {
	useSSL := false
	if env.MinioUseSSL == "true" {
		useSSL = true
	}

	// Inisialisasi MinIO client object
	minioClient, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
		Secure: useSSL,
	})
	
	if err != nil {
		log.Fatalf("Gagal terhubung ke MinIO: %v", err)
	}

	log.Println("MinIO terhubung dengan sukses!")
	return minioClient
}