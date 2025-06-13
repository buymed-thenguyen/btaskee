package db

import "time"

type Session struct {
	ID        uint `gorm:"primaryKey"`
	QuizID    uint
	Code      string
	CreatedAt time.Time
	StartAt   *time.Time
	CreatedBy uint
}

type Participant struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	QuizID       uint
	SessionID    uint
	TotalScore   int
	User         User `gorm:"foreignKey:UserID"`
	Quiz         Quiz `gorm:"foreignKey:QuizID"`
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DoneAt       *time.Time
	TimeConsumed int64 // millisecond
}

type ParticipantAnswer struct {
	ID             uint `gorm:"primaryKey"`
	ParticipantID  uint
	QuestionID     uint
	SessionID      uint
	AnswerOptionID uint
	IsCorrect      bool
	Participant    Participant `gorm:"foreignKey:ParticipantID"`
	Question       Question    `gorm:"foreignKey:QuestionID"`
	CreatedAt      time.Time
}
