package service

import (
	"errors"

	"github.com/duckcbuzz/crudapi/api/models"
	"github.com/jinzhu/gorm"
)

func ValidatePost(p *models.Post) error {

	if p.Title == "" {
		return errors.New("required title")
	}
	if p.Content == "" {
		return errors.New("required content")
	}
	if p.AuthorID < 1 {
		return errors.New("required author")
	}
	return nil
}

func SavePost(db *gorm.DB, p *models.Post) (*models.Post, error) {
	_, err := FindUserById(db, p.AuthorID)
	if err != nil {
		return &models.Post{}, err
	}
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &models.Post{}, err
	}
	return p, err
}

func FindAllPosts(db *gorm.DB) (*[]models.Post, error) {
	p := []models.Post{}
	err := db.Debug().Model(&models.Post{}).Preload("Author").Find(&p).Error
	if err != nil {
		return &[]models.Post{}, err
	}

	return &p, nil
}

func FindPostById(db *gorm.DB, pid uint32) (*models.Post, error) {
	p := models.Post{}
	err := db.Debug().Model(&models.Post{}).Preload("Author").Find(&p, pid).Error
	if err != nil {
		return &models.Post{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &models.Post{}, errors.New("user not found")
	}
	return &p, nil
}

func UpdatePost(db *gorm.DB, pid uint64, p *models.Post) (*models.Post, error) {
	p.ID = pid
	err := db.Debug().Save(p).Error
	if err != nil {
		return &models.Post{}, err
	}
	return p, nil
}

func DeletePost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&models.Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&models.Post{}).Delete(&models.Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
