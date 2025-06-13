package db

import (
	dbModel "btaskee/model/db"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InjectDB(db *gorm.DB) {
	DB = db
}

func GetUserByUsername(c *gin.Context, username string) (*dbModel.User, error) {
	var existing *dbModel.User
	err := DB.WithContext(c.Request.Context()).
		Where("username = ?", username).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func CreateUser(c *gin.Context, user *dbModel.User) error {
	return DB.WithContext(c.Request.Context()).Create(&user).Error
}

func GetUsersByIDs(c *gin.Context, ids []uint) ([]*dbModel.User, error) {
	var users []*dbModel.User
	err := DB.WithContext(c.Request.Context()).
		Where("id IN ?", ids).
		Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByID(c *gin.Context, id uint) (*dbModel.User, error) {
	var user *dbModel.User
	err := DB.WithContext(c.Request.Context()).
		Where("id = ?", id).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
