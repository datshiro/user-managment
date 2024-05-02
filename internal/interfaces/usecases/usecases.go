package usecases

import (
	"app/internal/interfaces/repositories"
	"app/internal/interfaces/usecases/user"

	"gorm.io/gorm"
)

type Usecases struct {
  UserUC user.UserUsecase 
}

func NewPostgresUsecase(dbc *gorm.DB) Usecases {
  return Usecases{
    UserUC: user.NewUseCase(repositories.NewRepo(dbc)),
  }
}
