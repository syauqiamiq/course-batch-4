package exercise

import (
	"course/internal/domain"
	"course/internal/dto"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseUsecase struct {
	db *gorm.DB
}

func NewExerciseUsecase(db *gorm.DB) *ExerciseUsecase {
	return &ExerciseUsecase{
		db: db,
	}
}

func (eu ExerciseUsecase) CreateNewExercise(c *gin.Context) {
	input := dto.ExerciseRequest{}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input",
		})
		return
	}

	exercices := domain.Exercise{
		Title:       input.Title,
		Description: input.Description,
	}

	if err := eu.db.Create(&exercices).Error; err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "failed creating the data",
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "Successfully create exercice",
	})
}

func (eu ExerciseUsecase) CreateNewQuestion(c *gin.Context) {
	stringID := c.Param("id")
	exerciseID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid param id",
		})
		return
	}

	input := dto.CreateQuestionRequest{}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input",
		})
		return
	}

	userID := c.Request.Context().Value("user_id").(int)

	question := domain.Question{
		ExerciseID:    exerciseID,
		Body:          input.Body,
		OptionA:       input.OptionA,
		OptionB:       input.OptionB,
		OptionC:       input.OptionC,
		OptionD:       input.OptionD,
		CorrectAnswer: input.CorrectAnswer,
		CreatorID:     userID,
	}

	if err := eu.db.Create(&question).Error; err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "failed creating the data",
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "Successfully create question",
	})
}

func (eu ExerciseUsecase) CreateNewAnswer(c *gin.Context) {
	input := dto.AnswerRequest{}

	exerciseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid param id",
		})
		return
	}
	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid param question id",
		})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input id",
		})
		return
	}

	userID := c.Request.Context().Value("user_id").(int)

	answer := domain.Answer{
		ExerciseID: exerciseID,
		QuestionID: questionID,
		UserID:     userID,
		Answer:     input.Answer,
	}

	if err := eu.db.Create(&answer).Error; err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "failed creating the data",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "Successfully create answer",
	})

}

func (eu ExerciseUsecase) GetExercise(c *gin.Context) {
	stringID := c.Param("id")
	exerciseID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input id",
		})
		return
	}

	var exercise domain.Exercise
	err = eu.db.Where("id = ?", exerciseID).Preload("Questions").Find(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "not found",
		})
		return
	}
	response := dto.CreateExerciseResponse(exercise)
	c.JSON(200, response)
}

func (eu ExerciseUsecase) GetScore(c *gin.Context) {
	stringID := c.Param("id")
	exerciseID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input id",
		})
		return
	}

	var exercise domain.Exercise
	err = eu.db.Where("id = ?", exerciseID).Preload("Questions").Find(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "not found",
		})
		return
	}

	userID := c.Request.Context().Value("user_id").(int)

	var answers []domain.Answer
	err = eu.db.Where("exercise_id = ? AND user_id = ?", exerciseID, userID).Find(&answers).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "not answered yet",
		})
		return
	}

	// calculate answer
	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score Score
	wg := new(sync.WaitGroup)
	for _, question := range exercise.Questions {
		wg.Add(1)
		go func(question domain.Question) {
			defer wg.Done()
			if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
				score.Inc(question.Score)
			}
		}(question)
	}
	wg.Wait()
	c.JSON(200, map[string]interface{}{
		"score": score.totalScore,
	})
}

type Score struct {
	totalScore int
	mu         sync.Mutex
}

func (s *Score) Inc(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalScore += value
}
