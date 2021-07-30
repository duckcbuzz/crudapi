package service

import (
	"github.com/duckcbuzz/crudapi/api/models"
	"github.com/jinzhu/gorm"
)

func SaveUser(db *gorm.DB, user *models.User) (*models.User, error) {
	err := db.Debug().Create(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}
