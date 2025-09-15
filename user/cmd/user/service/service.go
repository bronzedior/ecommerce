package service

import (
	"user/cmd/user/repository"
	"user/models"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (svc *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := svc.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *UserService) CreateNewUser(user *models.User) (int64, error) {
	userID, err := svc.UserRepo.InsertNewUser(user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
