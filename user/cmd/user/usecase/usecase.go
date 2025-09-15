package usecase

import (
	"context"
	"errors"
	"time"
	"user/cmd/user/service"
	"user/infrastructure/log"
	"user/models"
	"user/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type UserUsecase struct {
	UserService service.UserService
	JWTSecret   string
}

func NewUserUsecase(userService service.UserService, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		UserService: userService,
		JWTSecret:   jwtSecret,
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

func (uc *UserUsecase) Login(ctx context.Context, param *models.LoginParameter) (string, error) {
	user, err := uc.UserService.GetUserByEmail(ctx, param.Email)
	if user.ID == 0 {
		return "", errors.New("Email not found!")
	}
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("uc.UserService.GetUserByEmail got error: %v", err)
	}

	isMatch, err := utils.CheckPasswordHash(user.Password, param.Password)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("utils.CheckPasswordHash got error: %v", err)
	}

	if !isMatch {
		return "", errors.New("Email or password mismatched")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": param.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("utils.SignedString got error: %v", err)
	}

	return tokenString, nil
}
