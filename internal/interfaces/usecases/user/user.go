package user

import (
	"app/internal/interfaces/entities"
	"app/internal/interfaces/repositories"
	"app/internal/models"
	"context"
)

type UserUsecase interface {
	LoginWithEmail(ctx context.Context, email string, password string) (*models.User, error)
	LoginWithUsername(ctx context.Context, username string, password string) (*models.User, error)
	LoginWithPhone(ctx context.Context, phone string, password string) (*models.User, error)
	RegisterUser(ctx context.Context, data entities.UserData) (*models.User, error)
}

func NewUseCase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{}
}

type userUsecase struct {
	repo repositories.UserRepository
}

func (uc *userUsecase) LoginWithEmail(ctx context.Context, email string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) LoginWithUsername(ctx context.Context, username string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) LoginWithPhone(ctx context.Context, phone string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) RegisterUser(ctx context.Context, data entities.UserData) ( *models.User,  error) {
  userModel := data.ToModel()
  // TODO: Hash Password before store to db

  user , err := uc.repo.CreateUser(ctx, userModel)
  if err !=nil {
    return nil, err
  }
  return user, nil
}
