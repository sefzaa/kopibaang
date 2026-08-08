package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	
	// Database
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPassword             string `mapstructure:"DB_PASSWORD"`
	DBName                 string `mapstructure:"DB_NAME"`
	DBTimezone             string `mapstructure:"DB_TIMEZONE"`
	
	// JWT
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpireMins  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpireDays int    `mapstructure:"REFRESH_TOKEN_EXPIRE_DAYS"`
	
	// Redis
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                int    `mapstructure:"REDIS_DB"`
	
	// MinIO
	MinioEndpoint          string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey         string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey         string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucketName        string `mapstructure:"MINIO_BUCKET_NAME"`
	MinioUseSSL            string `mapstructure:"MINIO_USE_SSL"`

	// Firebase & SMTP
	FirebaseCredentials    string `mapstructure:"FIREBASE_CREDENTIALS"`
	SMTPEmail              string `mapstructure:"SMTP_EMAIL"`
	SMTPPassword           string `mapstructure:"SMTP_PASSWORD"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}