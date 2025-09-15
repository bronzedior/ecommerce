package repository

import (
	"errors"
	"user/models"

	"gorm.io/gorm"
)

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.Database.Where("email = ?", email).Last(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) InsertNewUser(user *models.User) (int64, error) {
	err := r.Database.Create(user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
