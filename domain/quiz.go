package domain

import (
	"btaskee/db"
	dbModel "btaskee/model/db"
	"btaskee/model/response"
	"btaskee/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func GetListQuiz(c *gin.Context) []*response.Quiz {
	quizzes, err := db.GetListQuiz(c)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return utils.MapSlice(quizzes, func(q *dbModel.Quiz) *response.Quiz {
		return &response.Quiz{
			ID:            q.ID,
			Title:         q.Title,
			TotalQuestion: len(q.Questions),
		}
	})
}

func GetQuizDetail(c *gin.Context, sessionCode string) *response.Quiz {
	session, err := db.GetSessionByCode(c, sessionCode)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if session == nil {
		c.Error(errors.New("session not found"))
		c.Status(http.StatusNotFound)
		return nil
	}
	if session.StartAt == nil || session.StartAt.After(time.Now()) {
		c.Error(errors.New("session not start yet"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	quiz, err := db.GetQuizByID(c, session.QuizID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if quiz == nil {
		c.Error(errors.New("quiz not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

	questions, err := db.GetQuestionsFromQuizID(c, session.QuizID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if len(questions) == 0 {
		c.Error(errors.New("questions not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

	for _, q := range questions {
		// Shuffle options
		rand.Shuffle(len(q.Options), func(i, j int) {
			q.Options[i], q.Options[j] = q.Options[j], q.Options[i]
		})
	}

	resp := response.Quiz{
		ID:    quiz.ID,
		Title: quiz.Title,
	}
	resp.Questions = utils.MapSlice(questions, func(q *dbModel.Question) *response.Question {
		return &response.Question{
			ID:           q.ID,
			QuestionText: q.QuestionText,
			AnswerOptions: utils.MapSlice(q.Options, func(a *dbModel.AnswerOption) *response.AnswerOption {
				return &response.AnswerOption{
					ID:   a.ID,
					Text: a.Text,
				}
			}),
		}
	})
	return &resp
}
