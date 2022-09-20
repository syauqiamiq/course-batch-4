package user

import (
	"course/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserUsecase struct {
	db *gorm.DB
}

func NewUserUsecase(db *gorm.DB) *UserUsecase {
	return &UserUsecase{
		db: db,
	}
}

func (uu UserUsecase) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid input",
		})
		return
	}

	if user.Name == "" {
		c.JSON(400, map[string]interface{}{
			"message": "name cannot be empty",
		})
		return
	}

	if user.Email == "" {
		c.JSON(400, map[string]interface{}{
			"message": "email cannot be empty",
		})
		return
	}

	if user.Password == "" {
		c.JSON(400, map[string]interface{}{
			"message": "password cannot be empty",
		})
		return
	}

	if len(user.Password) < 6 {
		c.JSON(400, map[string]interface{}{
			"message": "password must more thna 6 character",
		})
		return
	}

	if err := user.CreatePassword(user.Password); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "error when create user",
		})
		return
	}

	if err := uu.db.Create(&user).Error; err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error when create user",
		})
		return
	}
	token, err := user.GenerateToken()
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error when generate user token",
		})
		return
	}
	c.JSON(201, map[string]interface{}{
		"token": token,
	})
}

func (uu UserUsecase) Login(c *gin.Context) {
	var currentUser domain.User
	err := c.ShouldBind(&currentUser)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	var user domain.User
	err = uu.db.Where("email = ?", currentUser.Email).Take(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid email/password",
		})
		return
	}

	if err := user.ComparePassword(currentUser.Password); err != nil {
		c.JSON(400, gin.H{
			"message": "invalid email/password",
		})
		return
	}
	token, err := user.GenerateToken()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when get user",
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
}
