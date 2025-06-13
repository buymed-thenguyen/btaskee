package response

import "time"

type Participant struct {
	UserID       uint       `json:"user_id"`
	QuizID       uint       `json:"quiz_id"`
	UserName     string     `json:"user_name"`
	QuizName     string     `json:"quiz_name"`
	TotalScore   int        `json:"score"`
	DoneAt       *time.Time `json:"done_at"`
	SessionID    uint       `json:"session_id"`
	TimeConsumed int64      `json:"time_consumed"`
}

type Session struct {
	ID        uint       `json:"id"`
	QuizID    uint       `json:"quiz_id"`
	CreatedAt time.Time  `json:"created_at"`
	Code      string     `json:"code"`
	StartAt   *time.Time `json:"start_at"`
	CreatedBy uint       `json:"created_by"`
	Quiz      *Quiz      `json:"quiz"`
}

type ParticipantAnswer struct {
	ID             uint      `json:"id"`
	QuestionID     uint      `json:"question_id"`
	AnswerOptionID uint      `json:"answer_option_id"`
	CreatedAt      time.Time `json:"created_at"`
}
