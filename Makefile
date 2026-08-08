# Baca file .env
include .env
export

.PHONY: docker-dev docker-build docker-down docker-logs migrate-up migrate-down swagger

# 1. Perintah Harian (Cepat & Pendek)
docker-dev:
	@echo "Starting development containers..."
	@docker-compose up -d
	@echo "Development containers started"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"
	@echo "phpMyAdmin (MySQL): http://localhost:5050"
	@echo "Redis Commander: http://localhost:8081"
	@echo "MinIO Console: http://localhost:9001"

# 2. Perintah Rebuild (Dipakai HANYA jika ada perubahan di go.mod atau Dockerfile)
docker-build:
	@echo "Rebuilding and starting containers..."
	@docker-compose up --build -d
	@echo "Development containers started"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"
	@echo "phpMyAdmin (MySQL): http://localhost:5050"
	@echo "Redis Commander: http://localhost:8081"
	@echo "MinIO Console: http://localhost:9001"


# Mematikan seluruh container
docker-down:
	@echo "Stopping and removing containers..."
	@docker-compose down
	@echo "Containers stopped."

# Melihat logs dari aplikasi
docker-logs:
	docker-compose logs -f app

# Menjalankan migrasi database di dalam Docker
migrate-up:
	docker-compose exec app migrate -path database/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(mysql:3306)/${DB_NAME}" -verbose up

# Menurunkan migrasi database di dalam Docker
migrate-down:
	docker-compose exec app migrate -path database/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(mysql:3306)/${DB_NAME}" -verbose down

# Mengenerate dokumentasi Swagger
swagger:
	swag init -g cmd/main.go