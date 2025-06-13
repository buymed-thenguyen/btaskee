package db

import (
	dbModel "btaskee/model/db"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSessionByCode(c *gin.Context, code string) (*dbModel.Session, error) {
	var existing *dbModel.Session
	err := DB.WithContext(c.Request.Context()).
		Where("code = ?", code).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func CreateSession(c *gin.Context, session *dbModel.Session) error {
	return DB.WithContext(c.Request.Context()).Create(&session).Error
}

func UpdateSession(c *gin.Context, session *dbModel.Session) error {
	return DB.WithContext(c.Request.Context()).Save(session).Error
}
