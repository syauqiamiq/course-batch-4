package dto

import "course/internal/domain"

type ExerciseResponseDTO struct {
	ID          int                     `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Questions   QuestionListResponseDTO `json:"questions"`
}

func CreateExerciseResponse(e domain.Exercise) ExerciseResponseDTO {
	questionResp := CreateQuestionListResponseDTO(e.Questions)
	return ExerciseResponseDTO{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Questions:   questionResp,
	}
}

type QuestionResponseDTO struct {
	ID         int    `json:"id"`
	ExerciseID int    `json:"exercise_id"`
	Body       string `json:"body"`
	OptionA    string `json:"option_a"`
	OptionB    string `json:"option_b"`
	OptionC    string `json:"option_c"`
	OptionD    string `json:"option_d"`
	Score      int    `json:"score"`
	CreatorID  int    `json:"creator_id"`
}

func CreateQuestionResponseDTO(e domain.Question) QuestionResponseDTO {
	return QuestionResponseDTO{
		ID:         e.ID,
		ExerciseID: e.ExerciseID,
		Body:       e.Body,
		OptionA:    e.OptionA,
		OptionB:    e.OptionB,
		OptionC:    e.OptionC,
		OptionD:    e.OptionD,
		Score:      e.Score,
		CreatorID:  e.CreatorID,
	}
}

type QuestionListResponseDTO []QuestionResponseDTO

func CreateQuestionListResponseDTO(e []domain.Question) QuestionListResponseDTO {
	questionsResp := QuestionListResponseDTO{}
	for _, p := range e {
		question := CreateQuestionResponseDTO(p)
		questionsResp = append(questionsResp, question)
	}
	return questionsResp
}
