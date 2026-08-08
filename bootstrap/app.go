package bootstrap

import (
	"gorm.io/gorm"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Env   *Env
	DB    *gorm.DB
	Redis *redis.Client
	Minio *minio.Client // Sekarang MinIO sudah aktif
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewMySQLDatabase(app.Env)
	app.Redis = NewRedisClient(app.Env)
	
	// Tanda komentar dihilangkan, memanggil fungsi dari minio.go
	app.Minio = NewMinioClient(app.Env) 
	
	return *app
}

func (app *Application) CloseDBConnection() {
	if app.DB != nil {
		sqlDB, _ := app.DB.DB()
		sqlDB.Close()
	}
	if app.Redis != nil {
		app.Redis.Close()
	}
}