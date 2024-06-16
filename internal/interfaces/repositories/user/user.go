package user

import (
	"app/internal/interfaces/repositories/models"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

func NewRepo(dbc *gorm.DB, rd *redis.Client) UserRepository {
	return &userRepo{dbc, rd}
}

type userRepo struct {
	dbc *gorm.DB
	rd  *redis.Client
}

func (repo *userRepo) SaveUser(ctx context.Context, data *models.User) error {
	// remove from cache
	if err := repo.deleteFromCache(ctx, *data); err != nil {
		return err
	}

	// save user
	result := repo.dbc.Save(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *userRepo) GetUser(ctx context.Context, id int) (*models.User, error) {
	// get from cache
	cachedUser := repo.getUserFromCache(ctx, id)
	if cachedUser != nil {
		return cachedUser, nil
	}

	// Get user from database
	user := models.User{Model: gorm.Model{ID: uint(id)}}
	result := repo.dbc.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// cache retrieved user
	if err := repo.cacheUser(user); err != nil {
		return &user, err
	}
	return &user, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Get User from cache
	cachedUser := repo.getUserFromCacheByEmail(ctx, email)
	if cachedUser != nil {
		return cachedUser, nil
	}

	// Get from database
	user := models.User{Email: email}
	result := repo.dbc.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// cache retreived user
	if err := repo.cacheUser(user); err != nil {
		return &user, err
	}
	return &user, nil
}

func (repo *userRepo) DeleteUser(ctx context.Context, id int) error {
	// Delete from cache
	cachedUser := repo.getUserFromCache(ctx, id)
	if cachedUser != nil {
		if err := repo.deleteFromCache(ctx, *cachedUser); err != nil {
			return err
		}
	}

	result := repo.dbc.Delete(&models.User{}, id)
	return result.Error
}
