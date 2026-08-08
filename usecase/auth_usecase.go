package usecase

import (
	"context"
	"errors"
	"time"

	"kopibang/bootstrap"
	"kopibang/domain"
	"kopibang/domain/dto"
	"kopibang/domain/entity"
	"kopibang/internal/emailutils"
	"kopibang/internal/tokenutil"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo  domain.UserRepository
	redisRepo domain.RedisRepository
	env       *bootstrap.Env
}

func NewAuthUsecase(userRepo domain.UserRepository, redisRepo domain.RedisRepository, env *bootstrap.Env) *AuthUsecase {
	return &AuthUsecase{userRepo, redisRepo, env}
}

func (u *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest, expectedRole string) (dto.TokenResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.TokenResponse{}, errors.New("invalid credentials")
	}

	if user.Role != expectedRole {
		return dto.TokenResponse{}, errors.New("unauthorized access for this role")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return dto.TokenResponse{}, errors.New("invalid credentials")
	}

	accessToken, refreshToken, tokenID, err := tokenutil.GenerateTokens(user.ID, user.Role)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Simpan ke Redis (Expire 60 hari)
	err = u.redisRepo.SaveRefreshToken(ctx, user.ID, tokenID, 60*24*time.Hour)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         user.Role,
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	claims, err := tokenutil.ValidateToken(req.RefreshToken)
	if err != nil {
		return dto.TokenResponse{}, errors.New("invalid refresh token")
	}

	// Validasi token aktif di Redis
	err = u.redisRepo.ValidateRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return dto.TokenResponse{}, errors.New("refresh token has been revoked or expired")
	}

	// Hapus token lama di Redis untuk rotasi keamanan (One-time use)
	_ = u.redisRepo.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID)

	// Generate baru
	newAccess, newRefresh, newTokenID, err := tokenutil.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Simpan token baru
	err = u.redisRepo.SaveRefreshToken(ctx, claims.UserID, newTokenID, 60*24*time.Hour)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		Role:         claims.Role,
	}, nil
}

func (u *AuthUsecase) Register(ctx context.Context, req dto.RegisterRequest) (dto.TokenResponse, error) {
	// Cek apakah email sudah terdaftar
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return dto.TokenResponse{}, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	newUser := &entity.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "customer", // Hardcode karena ini pendaftaran user biasa
	}

	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Auto-Login: Generate token setelah berhasil insert
	accessToken, refreshToken, tokenID, err := tokenutil.GenerateTokens(newUser.ID, newUser.Role)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	err = u.redisRepo.SaveRefreshToken(ctx, newUser.ID, tokenID, 60*24*time.Hour)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         newUser.Role,
	}, nil
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	// 1. Pastikan email ada di DB
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("email not found")
	}

	// 2. Generate OTP
	otp := emailutils.GenerateOTP()

	// 3. Simpan di Redis (Expire 5 Menit)
	err = u.redisRepo.SaveOTP(ctx, req.Email, otp, 5*time.Minute)
	if err != nil {
		return errors.New("failed to process OTP")
	}

	// 4. Kirim Email pakai kredensial dari env
	go func() { 
		_ = emailutils.SendOTPResetPassword(req.Email, otp, u.env.SMTPEmail, u.env.SMTPPassword)
	}()

	return nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// 1. Validasi Reset Token (Bukan OTP lagi)
	err := u.redisRepo.ValidateResetToken(ctx, req.Email, req.ResetToken)
	if err != nil {
		return err
	}

	// 2. Cari User
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// 3. Hash Password Baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. Update Password di Database
	err = u.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return err
	}

	// 5. Hapus Reset Token setelah berhasil dipakai agar tidak bisa digunakan dua kali
	_ = u.redisRepo.DeleteResetToken(ctx, req.Email)

	return nil
}

func (u *AuthUsecase) Logout(ctx context.Context, req dto.LogoutRequest) error {
	// 1. Validasi token (meskipun mungkin sudah kedaluwarsa secara akses, kita butuh claims-nya)
	claims, err := tokenutil.ValidateToken(req.RefreshToken)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// 2. Hapus Refresh Token dari Redis
	err = u.redisRepo.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return errors.New("failed to logout, please try again")
	}

	return nil
}

func (u *AuthUsecase) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (dto.VerifyOTPResponse, error) {
	// 1. Validasi OTP di Redis
	err := u.redisRepo.ValidateOTP(ctx, req.Email, req.OTP)
	if err != nil {
		return dto.VerifyOTPResponse{}, err
	}

	// 2. Hapus OTP karena sudah terpakai dengan benar
	_ = u.redisRepo.DeleteOTP(ctx, req.Email)

	// 3. Generate Reset Token (berlaku 15 menit) agar user bisa ganti password
	resetToken := uuid.New().String()
	err = u.redisRepo.SaveResetToken(ctx, req.Email, resetToken, 15*time.Minute)
	if err != nil {
		return dto.VerifyOTPResponse{}, errors.New("failed to generate reset session")
	}

	return dto.VerifyOTPResponse{
		ResetToken: resetToken,
	}, nil
}