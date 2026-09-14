package server

import (
	"net/http"

	"sebosun/acrevus-go/internal/database"

	"github.com/gin-gonic/gin"
)

// func UserLogin()  {}
// func UserUpdate() {}
// func UserCreate() {}
// func UserDelete() {}

type UserForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (config *ApiConfig) UserCreate(c *gin.Context) {
	var userForm UserForm

	err := c.Copy().ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing payload"})
		return
	}

	hashedPassword, err := HashPassword(userForm.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong with the password"})
		return
	}

	userParams := database.CreateUserParams{
		Name:     userForm.Name,
		Email:    userForm.Email,
		Password: hashedPassword,
	}

	user, err := config.DB.CreateUser(c, userParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
