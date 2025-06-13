package domain

import (
	"btaskee/db"
	"btaskee/model/constant"
	"btaskee/model/request"
	"btaskee/model/response"
	"btaskee/utils"
	"btaskee/ws"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"time"

	dbModel "btaskee/model/db"
)

func CreateSessionWithQuizID(c *gin.Context, quizID uint) *response.Session {
	userID := c.GetUint(constant.USER_ID_CTX) // get user_id from jwt
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}

	session := &dbModel.Session{
		QuizID:    quizID,
		Code:      generateSessionCode(constant.SESSION_CODE_LEN),
		CreatedBy: userID,
	}
	if err := db.CreateSession(c, session); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	participant := &dbModel.Participant{
		UserID:    userID,
		QuizID:    quizID,
		SessionID: session.ID,
	}
	if err := db.CreateParticipant(c, participant); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return &response.Session{
		ID:        session.ID,
		QuizID:    quizID,
		Code:      session.Code,
		CreatedBy: userID,
	}
}

func generateSessionCode(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return strings.ToUpper(string(b))
}

func JoinSessionByCode(c *gin.Context, sessionCode string) *response.Participant {
	userID := c.GetUint(constant.USER_ID_CTX) // get user_id from jwt
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}
	user, err := db.GetUsersByID(c, userID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if user == nil {
		c.Error(errors.New("user not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

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
	if session.StartAt != nil && session.StartAt.Before(time.Now()) {
		c.Error(errors.New("session already started"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	// Check if participant already exists
	participant, err := db.GetParticipantByUserIDSessionID(c, userID, session.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if participant != nil {
		c.Error(errors.New("user already joined"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	participant = &dbModel.Participant{
		UserID:    userID,
		QuizID:    session.QuizID,
		SessionID: session.ID,
	}
	if err = db.CreateParticipant(c, participant); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	ws.BroadcastToSession(session.Code, "user_joined", gin.H{
		"user_name": user.Name,
	})

	return &response.Participant{
		UserID:    userID,
		QuizID:    session.QuizID,
		SessionID: session.ID,
		UserName:  user.Name,
	}
}

func GetLeaderboardBySession(c *gin.Context, sessionCode string) []*response.Participant {
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

	rows, err := db.GetSessionLeaderboard(c, session.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	userIDs := utils.MapSlice(rows, func(p *dbModel.Participant) uint {
		return p.UserID
	})
	users, err := db.GetUsersByIDs(c, userIDs)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	userMap := utils.SliceToMap(users, func(u *dbModel.User) uint {
		return u.ID
	})

	return utils.MapSlice(rows, func(r *dbModel.Participant) *response.Participant {
		var name string
		if user, exists := userMap[r.UserID]; exists {
			name = user.Name
		}
		return &response.Participant{
			UserID:       r.UserID,
			QuizID:       r.QuizID,
			SessionID:    r.SessionID,
			TotalScore:   r.TotalScore,
			DoneAt:       r.DoneAt,
			TimeConsumed: r.TimeConsumed,
			UserName:     name,
		}
	})
}

func SubmitAnswer(c *gin.Context, sessionCode string, req *request.SubmitAnswer) *response.DefaultResponse {
	userID := c.GetUint(constant.USER_ID_CTX) // get user_id from jwt
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}

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

	participant, err := db.GetParticipantByUserIDSessionID(c, userID, session.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if participant == nil {
		c.Error(errors.New("participant not found"))
		c.Status(http.StatusNotFound)
		return nil
	}
	if participant.DoneAt != nil && !participant.DoneAt.IsZero() {
		c.Error(errors.New("participant already done"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	totalScore, correctAnswerMap, err := calculateScore(c, session, req.Answers)
	if err != nil {
		return nil
	}

	participant.DoneAt = utils.ToPointerTime(time.Now())
	participant.TotalScore = totalScore
	participant.TimeConsumed = participant.DoneAt.Sub(*participant.CreatedAt).Milliseconds()

	err = db.UpdateParticipant(c, participant)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	participantAnswers := utils.MapSlice(req.Answers, func(a *request.Answer) *dbModel.ParticipantAnswer {
		isCorrect := correctAnswerMap[a.QuestionId].ID == a.AnswerOptionId
		return &dbModel.ParticipantAnswer{
			ParticipantID:  participant.ID,
			QuestionID:     a.QuestionId,
			AnswerOptionID: a.AnswerOptionId,
			SessionID:      session.ID,
			IsCorrect:      isCorrect,
		}
	})
	err = db.CreateParticipantAnswers(c, participantAnswers)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	ws.BroadcastToSession(session.Code, "user_answered", gin.H{
		"user_id":      userID,
		"session_code": sessionCode,
	})

	return &response.DefaultResponse{
		Message: "ok",
	}
}

func calculateScore(c *gin.Context, session *dbModel.Session, answers []*request.Answer) (int, map[uint]*dbModel.AnswerOption, error) {
	questions, err := db.GetQuestionsFromQuizID(c, session.QuizID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return 0, nil, nil
	}
	questionMap := utils.SliceToMap(questions, func(q *dbModel.Question) uint {
		return q.ID
	})

	questionIDs := utils.MapSlice(questions, func(q *dbModel.Question) uint {
		return q.ID
	})
	correctAnswers, err := db.GetCorrectAnswerByQuestionIDs(c, questionIDs)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return 0, nil, nil
	}

	correctAnswersMap := utils.SliceToMap(correctAnswers, func(a *dbModel.AnswerOption) uint {
		return a.QuestionID
	})

	var totalScore int
	for _, a := range answers {
		if a == nil {
			continue
		}
		question, exists := questionMap[a.QuestionId]
		if !exists {
			continue
		}
		if correctAnswersMap[a.QuestionId].ID == a.AnswerOptionId {
			totalScore += question.Score
		}
	}

	return totalScore, correctAnswersMap, nil
}

func StartSession(c *gin.Context, sessionCode string) *response.DefaultResponse {
	userID := c.GetUint(constant.USER_ID_CTX) // get user_id from jwt
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}

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
	if session.CreatedBy != userID {
		c.Error(errors.New("session not created by user"))
		c.Status(http.StatusBadRequest)
		return nil
	}
	if session.StartAt != nil && session.StartAt.Before(time.Now()) {
		c.Error(errors.New("session already started"))
		c.Status(http.StatusBadRequest)
		return nil
	}

	session.StartAt = utils.ToPointerTime(time.Now())
	if err = db.UpdateSession(c, session); err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	ws.BroadcastToSession(session.Code, "start_session", gin.H{
		"quiz_id": session.QuizID,
	})

	return &response.DefaultResponse{
		Message: "ok",
	}
}

func GetSessionDetail(c *gin.Context, sessionCode string) *response.Session {
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

	return &response.Session{
		ID:        session.ID,
		QuizID:    session.QuizID,
		Code:      sessionCode,
		CreatedBy: session.CreatedBy,
		Quiz: &response.Quiz{
			ID:            quiz.ID,
			Title:         quiz.Title,
			TotalQuestion: len(quiz.Questions),
		},
	}
}

func GetSessionParticipants(c *gin.Context, sessionCode string) []*response.Participant {
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

	participants, err := db.GetParticipantBySessionID(c, session.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	userIDs := utils.MapSlice(participants, func(p *dbModel.Participant) uint {
		return p.UserID
	})
	users, err := db.GetUsersByIDs(c, userIDs)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	userMap := utils.SliceToMap(users, func(u *dbModel.User) uint {
		return u.ID
	})

	return utils.MapSlice(participants, func(p *dbModel.Participant) *response.Participant {
		var userName string
		if user, exists := userMap[p.UserID]; exists {
			userName = user.Name
		}
		return &response.Participant{
			QuizID:       p.QuizID,
			SessionID:    p.SessionID,
			UserID:       p.UserID,
			UserName:     userName,
			TotalScore:   p.TotalScore,
			TimeConsumed: p.TimeConsumed,
		}
	})
}

func GetSessionParticipantAnswers(c *gin.Context, sessionCode string) []*response.ParticipantAnswer {
	userID := c.GetUint(constant.USER_ID_CTX) // get user_id from jwt
	if userID == 0 {
		c.Error(errors.New("missing user_id in context"))
		c.Status(http.StatusUnauthorized)
		return nil
	}

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

	participant, err := db.GetParticipantByUserIDSessionID(c, userID, session.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}
	if participant == nil {
		c.Error(errors.New("participant not found"))
		c.Status(http.StatusNotFound)
		return nil
	}

	participantAnswers, err := db.GetParticipantAnswersByParticipantID(c, participant.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return nil
	}

	return utils.MapSlice(participantAnswers, func(a *dbModel.ParticipantAnswer) *response.ParticipantAnswer {
		return &response.ParticipantAnswer{
			ID:             a.ID,
			QuestionID:     a.QuestionID,
			AnswerOptionID: a.AnswerOptionID,
			CreatedAt:      a.CreatedAt,
		}
	})
}
