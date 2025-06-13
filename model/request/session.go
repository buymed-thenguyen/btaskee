package request

type Answer struct {
	QuestionId     uint `json:"question_id"`
	AnswerOptionId uint `json:"answer_option_id"`
}

type SubmitAnswer struct {
	Answers []*Answer `json:"answers"`
}

type CreateSession struct {
	QuizID uint `json:"quiz_id"`
}
