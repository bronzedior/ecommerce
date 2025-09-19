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

func (uc *UserUsecase) GetUserByUserID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := uc.UserService.GetUserByUserID(ctx, userID)
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

func (uc *UserUsecase) Login(ctx context.Context, param models.LoginParameter, userID int64, storedPassword string) (string, error) {
	isMatch, err := utils.CheckPasswordHash(storedPassword, param.Password)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("utils.CheckPasswordHash got error: %v", err)
	}

	if !isMatch {
		return "", errors.New("email or password mismatched")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": param.Email,
		}).Errorf("token.SignedString got error: %v", err)
		return "", err
	}

	return tokenString, nil
}
