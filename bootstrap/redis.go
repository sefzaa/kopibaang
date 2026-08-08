package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env *Env) *redis.Client {
	// Format host dan port secara dinamis dari .env
	addr := fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: env.RedisPassword,
		DB:       env.RedisDB, 
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Gagal terhubung ke Redis di %s: %v", addr, err)
	}

	log.Println("Redis terhubung dengan sukses!")
	return rdb
}