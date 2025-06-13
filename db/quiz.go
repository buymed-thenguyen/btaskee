package db

import (
	dbModel "btaskee/model/db"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCorrectAnswerByQuestionIDs(c *gin.Context, questionIDs []uint) ([]*dbModel.AnswerOption, error) {
	var rows []*dbModel.AnswerOption
	err := DB.WithContext(c.Request.Context()).
		Where("question_id IN ?", questionIDs).
		Where("is_correct = ?", true).
		Find(&rows).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetQuestionsFromQuizID(c *gin.Context, quizID uint) ([]*dbModel.Question, error) {
	var rows []*dbModel.Question
	err := DB.WithContext(c.Request.Context()).
		Where("quiz_id = ?", quizID).
		Preload("Options").
		Find(&rows).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetQuizByID(c *gin.Context, quizID uint) (*dbModel.Quiz, error) {
	var quiz *dbModel.Quiz
	err := DB.WithContext(c.Request.Context()).
		Where("id = ?", quizID).
		Preload("Questions").
		First(&quiz).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func GetListQuiz(c *gin.Context) ([]*dbModel.Quiz, error) {
	var quiz []*dbModel.Quiz
	err := DB.WithContext(c.Request.Context()).
		Order("id desc").
		Preload("Questions").
		Find(&quiz).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return quiz, nil
}
