package repositories

import (
	"app/internal/interfaces/entities"
	"app/internal/models"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
  CreateUser(ctx context.Context, data entities.UserData) (*models.User, error)
  GetUserByID(ctx context.Context, id int) (*models.User, error)
  UpdateLastLogin(ctx context.Context, id int) (*models.User, error)
}

func NewRepo(dbc *gorm.DB) UserRepository {
  return &userRepo{dbc}
}

type userRepo struct { 
  dbc *gorm.DB
}

func (repo *userRepo) CreateUser(ctx context.Context, data entities.UserData) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (repo *userRepo) GetUserByID(ctx context.Context, id int) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (repo *userRepo) UpdateLastLogin(ctx context.Context, id int) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

