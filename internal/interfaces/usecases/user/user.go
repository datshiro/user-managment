package user

import (
	"app/internal/interfaces/repositories/models"
	"app/internal/interfaces/repositories/user"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidPassword = errors.New("login: password is incorrect")
  ErrUserExisted = errors.New("create: user is already existed")
)

type UserUsecase interface {
	LoginWithEmail(ctx context.Context, email string, password string) (*models.User, error)
	RegisterUser(ctx context.Context, data *models.User) (*models.User, error)
}

func NewUseCase(repo user.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

type userUsecase struct {
	repo user.UserRepository
}

func (uc *userUsecase) LoginWithEmail(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.ComparePassword(password) {
		return nil, ErrInvalidPassword
	}

	user.LatestLogin = time.Now()
	err = uc.repo.SaveUser(ctx, user)
	return user, err
}

func (uc *userUsecase) RegisterUser(ctx context.Context, data *models.User) (*models.User, error) {
	// Check if user existed
	u, err := uc.repo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
  if u != nil {
    return nil, ErrUserExisted
  }

	if err = uc.repo.SaveUser(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}
