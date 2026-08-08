package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"kopibang/domain"
	"kopibang/domain/dto"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo}
}

func (u *UserUsecase) GetProfile(ctx context.Context, userIDStr string) (*dto.UserProfileResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Points:   user.Points,
		Role:     user.Role,
	}, nil
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, userIDStr string, req dto.UpdateProfileRequest) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update data yang diizinkan
	user.Name = req.Name
	user.Username = req.Username

	return u.userRepo.Update(ctx, user)
}