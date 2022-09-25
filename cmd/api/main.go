package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"message": "hello world",
		})
	})
	db := database.CreateConn()
	exercise := exercise.NewExerciseUsecase(db)
	// exercise endpoint
	r.POST("/exercises", middleware.WithAuth(), exercise.CreateNewExercise)
	r.GET("/exercises/:id", middleware.WithAuth(), exercise.GetExercise)
	r.GET("/exercises/:id/scores", middleware.WithLog(), middleware.WithAuth(), exercise.GetScore)
	r.POST("/exercises/:id/questions", middleware.WithAuth(), exercise.CreateNewQuestion)
	r.POST("/exercises/:id/questions/:questionId/answer", middleware.WithAuth(), exercise.CreateNewAnswer)

	// user endpoint
	userUsecase := user.NewUserUsecase(db)
	r.POST("/register", userUsecase.Register)
	r.Run(":1234")
}
