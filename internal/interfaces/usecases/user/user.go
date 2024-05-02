package user

import (
	"app/internal/interfaces/repositories"
	"app/internal/models"
	"log"
)

type UserUsecase interface{
  LoginWithEmail(email string, password string) (*models.User, error)
  LoginWithUsername(username string, password string) (*models.User, error)
  LoginWithPhone(phone string, password string) (*models.User, error)
  RegisterUser() (*models.User, error)
}

func NewUseCase(repo repositories.UserRepository) UserUsecase {
  return &userUsecase{}
}

type userUsecase struct {}

func (uc *userUsecase) LoginWithEmail(email string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) LoginWithUsername(username string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) LoginWithPhone(phone string, password string) (_ *models.User, _ error) {
	panic("not implemented") // TODO: Implement
}

func (uc *userUsecase) RegisterUser() (_ *models.User, _ error) {
  log.Println("RegisterUser")
  return nil,nil
}

