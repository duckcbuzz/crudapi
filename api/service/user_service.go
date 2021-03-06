package service

import (
	"encoding/base64"
	"errors"
	"log"
	"os"

	"github.com/badoux/checkmail"
	"github.com/duckcbuzz/crudapi/api/models"
	"github.com/jinzhu/gorm"
)

func ValidateUser(u *models.User) error {
	if u.Nickname == "" {
		return errors.New("required nickname")
	}
	if u.Password == "" {
		return errors.New("required password")
	}
	if u.Email == "" {
		return errors.New("iequired email")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("invalid email")
	}
	return nil
}

func BeforeSave(user *models.User) error {
	err := ValidateUser(user)
	if err != nil {
		return err
	}
	hashedPassword, err := models.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func SaveUser(db *gorm.DB, user *models.User) (*models.User, error) {
	err := BeforeSave(user)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Create(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func FindAllUsers(db *gorm.DB) (*[]models.User, error) {
	users := []models.User{}
	err := db.Debug().Model(&models.User{}).Find(&users).Error
	if err != nil {
		return &[]models.User{}, err
	}
	return &users, nil
}

func FindUserById(db *gorm.DB, uid uint32) (*models.User, error) {
	user := models.User{}
	err := db.Debug().Model(&models.User{}).First(&user, uid).Error
	if err != nil {
		return &models.User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &models.User{}, errors.New("user not found")
	}
	return &user, nil
}

func FindByUsername(db *gorm.DB, username string) (*models.User, error) {
	user := models.User{}
	err := db.Debug().Model(&models.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &models.User{}, errors.New("user not found")
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *models.User, uid uint32) (*models.User, error) {
	err := BeforeSave(user)
	if err != nil {
		log.Fatal(err)
	}
	user.ID = uid
	err = db.Debug().Save(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func DeleteUser(db *gorm.DB, user *models.User) (int64, error) {
	err := db.Debug().Delete(&user).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func DeleteUserById(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&models.User{}).Where("id = ?", uid).Delete(&models.User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func ResetPassword(db *gorm.DB, email string) error {
	user := models.User{}
	err := db.Debug().Model(&models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return errors.New("user not found")
	}

	s := "new password"
	se := base64.StdEncoding.EncodeToString([]byte(s))

	to := []string{user.Email}
	from := os.Getenv("MAIL_DOMAIN")
	msg := "From: " + from + "\n" +
		"To: " + to[0] + "\n" +
		"Subject:Reset password\n\n" +
		"New password: " + se
	err = SendMail(to, []byte(msg))
	if err != nil {
		return err
	}

	hashedPassword, err := models.Hash(se)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = db.Debug().Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}
