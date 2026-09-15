package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (config *APIConfig) AuthRequired(tokenString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, prefix)

		userID, err := config.verifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Make authenticated user available to later handlers
		c.Set("userID", userID)

		c.Next()
	}
}

func (config *APIConfig) verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(config.ENVS.JwtSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (config *APIConfig) signJWT(userID int) (string, error) {
	secret := []byte(config.ENVS.JwtSecret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}

	return s, nil
}
