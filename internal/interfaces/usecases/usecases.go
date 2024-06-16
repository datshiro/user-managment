package usecases

import (
	user_repo "app/internal/interfaces/repositories/user"
	"app/internal/interfaces/usecases/user"

)

type Usecases struct {
  UserUC user.UserUsecase 
}

func NewPostgresUsecase(userRepo user_repo.UserRepository) Usecases {
  return Usecases{
    UserUC: user.NewUseCase(userRepo),
  }
}
