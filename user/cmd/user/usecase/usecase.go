package usecase

import (
	"context"
	"user/cmd/user/service"
	"user/infrastructure/log"
	"user/models"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type UserUsecase struct {
	UserService service.UserService
}

func NewUserUsecase(userService service.UserService) *UserUsecase {
	return &UserUsecase{
		UserService: userService,
	}
}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := uc.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
		}).Errorf("utils.Hashpassword() got error %v", err)
		return err
	}

	user.Password = hashedPassword
	_, err = uc.UserService.CreateNewUser(ctx, user)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
			"name":  user.Name,
		}).Errorf("uc.UserService.CreateNewUser() got error %v", err)
		return err
	}

	return nil
}
