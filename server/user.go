package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"sebosun/acrevus-go/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UserForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateForm struct {
	Name     *string `json:"name" binding:"omitempty,min=1"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=1"`
}

// UserResponse is the public representation of a user, excluding the password hash.
type UserResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func userResponse(user database.User) UserResponse {
	response := UserResponse{
		ID: user.ID, Name: user.Name, Email: user.Email,
		CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	}
	if user.DeletedAt.Valid {
		response.DeletedAt = &user.DeletedAt.Time
	}
	return response
}

func (config *ApiConfig) registerUserRoutes(router *gin.RouterGroup) {
	router.POST("/users", config.UserCreate)
	router.GET("/users", config.UserList)
	router.GET("/users/:id", config.UserGet)
	router.PATCH("/users/:id", config.UserUpdate)
	router.DELETE("/users/:id", config.UserDelete)
}

func (config *ApiConfig) UserCreate(c *gin.Context) {
	var userForm UserForm

	err := c.ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, a valid email, and password are required"})
		return
	}

	hashedPassword, err := HashPassword(userForm.Password)
	if err != nil {
		userError(c, err)
		return
	}

	userParams := database.CreateUserParams{
		Name:     userForm.Name,
		Email:    userForm.Email,
		Password: hashedPassword,
	}

	user, err := config.DB.CreateUser(c.Request.Context(), userParams)

	if err != nil {
		userError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse(user)})
}

func (config *ApiConfig) UserGet(c *gin.Context) {
	id, ok := GetUserID(c)

	if !ok {
		return
	}

	user, err := config.DB.GetUser(c.Request.Context(), id)
	if err != nil {
		userError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse(user)})
}

func (config *ApiConfig) UserList(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be a non-negative integer"})
		return
	}
	users, err := config.DB.ListUsers(c.Request.Context(), database.ListUsersParams{
		Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		userError(c, err)
		return
	}
	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, userResponse(user))
	}
	c.JSON(http.StatusOK, gin.H{"users": response, "limit": limit, "offset": offset})
}

func (config *ApiConfig) UserUpdate(c *gin.Context) {
	id, ok := GetUserID(c)
	if !ok {
		return
	}
	var form UserUpdateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide non-empty fields and a valid email"})
		return
	}
	if form.Name == nil && form.Email == nil && form.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide at least one of name, email, or password"})
		return
	}
	params := database.UpdateUserParams{ID: id}
	if form.Name != nil {
		params.Name = sql.NullString{String: *form.Name, Valid: true}
	}
	if form.Email != nil {
		params.Email = sql.NullString{String: *form.Email, Valid: true}
	}
	if form.Password != nil {
		hash, err := HashPassword(*form.Password)
		if err != nil {
			userError(c, err)
			return
		}
		params.Password = sql.NullString{String: hash, Valid: true}
	}
	user, err := config.DB.UpdateUser(c.Request.Context(), params)
	if err != nil {
		userError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse(user)})
}

func (config *ApiConfig) UserDelete(c *gin.Context) {
	id, ok := GetUserID(c)
	if !ok {
		return
	}
	rows, err := config.DB.SoftDeleteUser(c.Request.Context(), id)
	if err != nil {
		userError(c, err)
		return
	}
	if rows == 0 {
		userError(c, sql.ErrNoRows)
		return
	}
	c.Status(http.StatusNoContent)
}

func GetUserID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID must be a positive integer"})
		return 0, false
	}

	return id, true
}

func userError(c *gin.Context, err error) {
	var pgErr *pq.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, ErrPasswordTooLong):
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrPasswordTooLong.Error()})
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already in use"})
	default:
		log.Printf("user request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process user request"})
	}
}
