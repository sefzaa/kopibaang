package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"kopibang/domain"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) SaveResetToken(ctx context.Context, email string, token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("reset_token:%s", email)
	return r.client.Set(ctx, key, token, expiresIn).Err()
}

func (r *redisRepository) ValidateResetToken(ctx context.Context, email string, token string) error {
	key := fmt.Sprintf("reset_token:%s", email)
	savedToken, err := r.client.Get(ctx, key).Result()
	if err != nil { return fmt.Errorf("reset token expired or invalid") }
	if savedToken != token { return fmt.Errorf("invalid reset token") }
	return nil
}

func (r *redisRepository) DeleteResetToken(ctx context.Context, email string) error {
	key := fmt.Sprintf("reset_token:%s", email)
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) SaveOTP(ctx context.Context, email string, otp string, expiresIn time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.client.Set(ctx, key, otp, expiresIn).Err()
}

func (r *redisRepository) ValidateOTP(ctx context.Context, email string, otp string) error {
	key := fmt.Sprintf("otp:%s", email)
	savedOTP, err := r.client.Get(ctx, key).Result()
	if err != nil { return fmt.Errorf("OTP expired or does not exist") }
	if savedOTP != otp { return fmt.Errorf("invalid OTP") }
	return nil
}

func (r *redisRepository) DeleteOTP(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) SetState(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err() 
}

func (r *redisRepository) GetState(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// --- TAMBAHAN UNTUK REFRESH TOKEN & QR TOKEN ---
func (r *redisRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return r.client.Set(ctx, key, "valid", expiresIn).Err()
}

func (r *redisRepository) ValidateRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return r.client.Get(ctx, key).Err()
}

func (r *redisRepository) DeleteRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) ValidateQRToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("qr:%s", token)
	return r.client.Get(ctx, key).Result()
}


func (r *redisRepository) SaveQRToken(ctx context.Context, tokenType string, token string, data string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", tokenType, token)
	return r.client.Set(ctx, key, data, expiresIn).Err()
}

func (r *redisRepository) GetQRTokenData(ctx context.Context, tokenType string, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", tokenType, token)
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepository) DeleteQRToken(ctx context.Context, tokenType string, token string) error {
	key := fmt.Sprintf("%s:%s", tokenType, token)
	return r.client.Del(ctx, key).Err()
}