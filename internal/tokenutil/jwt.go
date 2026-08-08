package tokenutil

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	SecretKey = []byte("super-secret-key-ganti-di-env") // Nanti idealnya ambil dari .env
)

type JwtCustomClaims struct {
	UserID  uuid.UUID `json:"user_id"`
	Role    string    `json:"role"`
	TokenID string    `json:"jti"` // JWT ID untuk validasi di Redis
	jwt.RegisteredClaims
}

func GenerateTokens(userID uuid.UUID, role string) (string, string, string, error) {
	tokenID := uuid.New().String()

	// 1. Access Token (Singkat: 15 menit)
	accessClaims := &JwtCustomClaims{
		UserID:  userID,
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(SecretKey)
	if err != nil {
		return "", "", "", err
	}

	// 2. Refresh Token (Panjang: 60 hari)
	refreshClaims := &JwtCustomClaims{
		UserID:  userID,
		Role:    role,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SecretKey)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, tokenID, nil
}

func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}