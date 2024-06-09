package user

import (
	"app/internal/interfaces/repositories"
	"app/internal/models"
	"context"
)

type UserUsecase interface {
	LoginWithEmail(ctx context.Context, email string, password string) (*models.User, error)
	LoginWithUsername(ctx context.Context, username string, password string) (*models.User, error)
	LoginWithPhone(ctx context.Context, phone string, password string) (*models.User, error)
	RegisterUser(ctx context.Context, data *models.User) (*models.User, error)
}

func NewUseCase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

type userUsecase struct {
	repo repositories.UserRepository
}

func (uc *userUsecase) LoginWithEmail(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email, password)
	if err != nil {
		return nil, err
	}
	_, _ = uc.repo.UpdateLastLogin(ctx, user.ID)
	return user, err
}

func (uc *userUsecase) LoginWithUsername(ctx context.Context, username string, password string) (_ *models.User, _ error) {
	user, err := uc.repo.GetUserByUsername(ctx, username, password)
	if err != nil {
		return nil, err
	}
	_, _ = uc.repo.UpdateLastLogin(ctx, user.ID)
	return user, err
}

func (uc *userUsecase) LoginWithPhone(ctx context.Context, phone string, password string) (_ *models.User, _ error) {
	user, err := uc.repo.GetUserByPhone(ctx, phone, password)
	if err != nil {
		return nil, err
	}
	_, _ = uc.repo.UpdateLastLogin(ctx, user.ID)
	return user, err
}

func (uc *userUsecase) RegisterUser(ctx context.Context, data *models.User) (*models.User, error) {
	// Hash Password implemented at BeforeCreate hook
	user, err := uc.repo.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return user, nil
}
