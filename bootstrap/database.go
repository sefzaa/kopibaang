package bootstrap

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDatabase(env *Env) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL via GORM:", err)
	}

	log.Println("Successfully connected to MySQL via GORM")
	return db
}