package response

type Quiz struct {
	ID            uint        `json:"id"`
	Title         string      `json:"title"`
	Questions     []*Question `json:"questions"`
	TotalQuestion int         `json:"total_question"`
}

type Question struct {
	ID            uint            `json:"id"`
	QuestionText  string          `json:"question_text"`
	AnswerOptions []*AnswerOption `json:"answer_options"`
}

type AnswerOption struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}
