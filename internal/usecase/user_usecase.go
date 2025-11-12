package usecase

import (
	"test-service/internal/domain"
	"test-service/pkg/logger"
)

type userUsecase struct {
	repo   domain.UserRepository
	logger logger.Logger
}

func NewUserUsecase(repo domain.UserRepository, logr logger.Logger) domain.UserUsecase {
	return &userUsecase{repo: repo, logger: logr}
}

func (u *userUsecase) CreateUser(name, email string) error {
	user := &domain.User{Name: name, Email: email}
	err := u.repo.Create(user)
	if err != nil {
		u.logger.Errorf("failed to create user: %v", err)
		return err
	}
	return nil
}

func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecase) GetUserByID(id int64) (*domain.User, error) {
	return u.repo.FindByID(id)
}
