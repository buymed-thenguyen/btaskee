package db

type Quiz struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Questions []*Question
}

type Question struct {
	ID           uint   `gorm:"primaryKey"`
	QuizID       uint   `gorm:"index"`
	QuestionText string `gorm:"not null"`
	Score        int
	Options      []*AnswerOption
}

type AnswerOption struct {
	ID         uint
	QuestionID uint
	Text       string
	IsCorrect  bool
}
