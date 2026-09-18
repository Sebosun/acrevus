package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"sebosun/acrevus-go/helpers"
	"sebosun/acrevus-go/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

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

func (config *APIConfig) FetchMe(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)

	if !ok {
		return
	}

	user, err := config.DB.GetUser(c.Request.Context(), userID)
	if err != nil {
		userError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse(user)})
}

func (config *APIConfig) UserGet(c *gin.Context) {
	id, ok := GetIDFromParams(c)

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

func (config *APIConfig) UserList(c *gin.Context) {
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

func (config *APIConfig) UserUpdate(c *gin.Context) {
	id, ok := GetIDFromParams(c)
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
		hash, err := helpers.HashPassword(*form.Password)
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

func (config *APIConfig) UserDelete(c *gin.Context) {
	id, ok := GetIDFromParams(c)
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

func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(KeysUserID)
	id, ok := userID.(int64)

	if !exists || !ok || id <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Something wrong with authorization"})
		return 0, false
	}

	return id, true
}

func GetIDFromParams(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is either missing or is not a number"})
		return 0, false
	}

	return id, true
}

func userError(c *gin.Context, err error) {
	var pgErr *pq.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.Is(err, helpers.ErrPasswordTooLong):
		c.JSON(http.StatusBadRequest, gin.H{"error": helpers.ErrPasswordTooLong.Error()})
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already in use"})
	default:
		log.Printf("user request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process user request"})
	}
}
