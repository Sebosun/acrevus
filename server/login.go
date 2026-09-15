package server

import (
	"fmt"
	"net/http"

	"sebosun/acrevus-go/helpers"

	"github.com/gin-gonic/gin"
)

type UserLoginForm struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=1"`
}

func (config *APIConfig) UserLogin(c *gin.Context) {
	var userForm UserLoginForm

	err := c.ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, a valid email, and password are required"})
		return
	}

	user, err := config.DB.GetUserByEmail(c.Request.Context(), *userForm.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	isPasswordGucci := helpers.CheckPassword(*userForm.Password, user.Password)
	if !isPasswordGucci {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	signedJWT, err := config.signJWT(int(user.ID))
	if err != nil {
		fmt.Printf("Error parsing signedJWT %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userResponse(user),
		"token": signedJWT,
	})
}
