package db

import (
	dbModel "btaskee/model/db"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetParticipantAnswersByParticipantID(c *gin.Context, participantID uint) ([]*dbModel.ParticipantAnswer, error) {
	var answers []*dbModel.ParticipantAnswer
	err := DB.WithContext(c.Request.Context()).
		Where("participant_id = ?", participantID).
		Find(&answers).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func CreateParticipantAnswers(c *gin.Context, answers []*dbModel.ParticipantAnswer) error {
	return DB.WithContext(c.Request.Context()).Create(&answers).Error
}
