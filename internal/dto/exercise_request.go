package dto

type ExerciseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AnswerRequest struct {
	Answer string `json:"answer"`
}
