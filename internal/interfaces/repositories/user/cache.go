package user

import (
	"app/internal/interfaces/repositories/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func getIdKey(id uint) string {
	return fmt.Sprintf("users:%d", id)
}
func getEmailKey(email string) string {
	return fmt.Sprintf("users:%s", email)
}

func (repo *userRepo) convertMapToUse(mapUser map[string]string) (models.User, error) {
	userID, err := strconv.Atoi(mapUser["id"])
	if err != nil {
		return models.User{}, err
	}

	createdAt, err := time.Parse(time.RFC3339Nano, mapUser["created_at"])
	if err != nil {
		return models.User{}, err
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, mapUser["updated_at"])
	if err != nil {
		return models.User{}, err
	}

	birthDay, err := time.Parse(time.RFC3339Nano, mapUser["birthday"])
	if err != nil {
		return models.User{}, err
	}

	latestLogin, err := time.Parse(time.RFC3339Nano, mapUser["latest_login"])
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Model: gorm.Model{
			ID:        uint(userID),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Fullname:    mapUser["fullname"],
		PhoneNumber: mapUser["phone_number"],
		Email:       mapUser["email"],
		Password:    mapUser["password"],
		Birthday:    birthDay,
		LatestLogin: latestLogin,
	}
	return user, nil
}

func (repo *userRepo) cacheUser(user models.User) error {
	idKey := getIdKey(user.ID)
	emailKey := getEmailKey(user.Email)
	cachedUser := map[string]any{
		"id":           user.ID,
		"email":        user.Email,
		"fullname":     user.Fullname,
		"phone_number": user.PhoneNumber,
		"password":     user.Password,
		"birthday":     user.Birthday,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
		"latest_login": user.LatestLogin,
	}

	for key, value := range cachedUser {
		if err := repo.rd.HSet(context.Background(), idKey, key, value).Err(); err != nil {
			return err
		}
		if err := repo.rd.HSet(context.Background(), emailKey, key, value).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (repo *userRepo) getUserFromCache(ctx context.Context, id int) *models.User {
	result := repo.rd.HGetAll(ctx, getIdKey(uint(id))).Val()
	if len(result) == 0 {
		return nil
	}
	user, err := repo.convertMapToUse(result)
	if err != nil {
		panic(err)
	}
	return &user
}

func (repo *userRepo) getUserFromCacheByEmail(ctx context.Context, email string) *models.User {
	result := repo.rd.HGetAll(ctx, getEmailKey(email)).Val()
	if len(result) == 0 {
		return nil
	}
	user, err := repo.convertMapToUse(result)
	if err != nil {
    panic(err)
	}
	return &user
}

func (repo *userRepo) deleteFromCache(ctx context.Context, user models.User) error {
	keysToDelete := make([]string, 0)
	if user.ID > 0 {
		keysToDelete = append(keysToDelete, getIdKey(user.ID))
	}
	if user.Email != "" {
		keysToDelete = append(keysToDelete, getEmailKey(user.Email))
	}
	return repo.rd.Del(ctx, keysToDelete...).Err()
}
