package repositories

import (
	"app/internal/consts"
	"app/internal/models"
	"app/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userObject *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id uint) (bool, error)
	GetUserByEmail(ctx context.Context, email string, password string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string, password string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string, password string) (*models.User, error)
}

func NewRepo(dbc *gorm.DB) UserRepository {
	return &userRepo{dbc}
}

type userRepo struct {
	dbc *gorm.DB
}

func (repo *userRepo) CreateUser(ctx context.Context, data *models.User) (*models.User, error) {
	tx := repo.dbc.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check user existence
	var existed bool
	err := tx.Model(&models.User{}).
		Select("count(*) > 0").
		Where("username = ? OR email = ? OR phone_number = ?", data.Username, data.Email, data.PhoneNumber).
		Find(&existed).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consts.NewError(fmt.Errorf("user already existed")).WithTag("Method", "CreateUser")
		} else {
			return nil, consts.NewError(fmt.Errorf("something went wrong. Please try again; %+v", err)).WithTag("Method", "Commit")
		}
	}
	if existed {
		return nil, consts.NewError(fmt.Errorf("user already existed")).WithTag("Method", "CreateUser")
	}

	result := tx.Create(data)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if err := tx.Commit().Error; err != nil {
		return nil, consts.NewError(fmt.Errorf("something went wrong. Please try again; %+v", err)).WithTag("Method", "Commit")
	}
	return data, nil
}

func (repo *userRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	result := repo.dbc.First(user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *userRepo) UpdateLastLogin(ctx context.Context, id uint) (bool, error) {
	result := repo.dbc.Model(&models.User{}).Where("id = ?", id).Update("latest_login", time.Now())
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string, password string) (*models.User, error) {
	user := &models.User{}
	result := repo.dbc.First(user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	if ok := utils.CheckPasswordHash(password, user.Password); !ok {
		return nil, consts.NewError(fmt.Errorf("credential verification failure; wrong password")).WithTag("Method", "GetUserByEmail")
	}
	return user, nil
}

func (repo *userRepo) GetUserByPhone(ctx context.Context, phone string, password string) (*models.User, error) {
	user := &models.User{}
	result := repo.dbc.First(user, "phone_number = ?", phone)
	if result.Error != nil {
		return nil, result.Error
	}
	if ok := utils.CheckPasswordHash(password, user.Password); !ok {
		return nil, consts.NewError(fmt.Errorf("credential verification failure; wrong password")).WithTag("Method", "GetUserByEmail")
	}
	return user, nil
}

func (repo *userRepo) GetUserByUsername(ctx context.Context, username string, password string) (*models.User, error) {
	user := &models.User{}
	result := repo.dbc.First(user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	if ok := utils.CheckPasswordHash(password, user.Password); !ok {
		return nil, consts.NewError(fmt.Errorf("credential verification failure; wrong password")).WithTag("Method", "GetUserByEmail")
	}
	return user, nil
}
