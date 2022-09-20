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
