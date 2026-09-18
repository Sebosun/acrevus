package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	Secret string
}

type Claims struct {
	UserID int64 `json:"user_id"`

	jwt.RegisteredClaims
}

func (jwtService *JwtService) VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(jwtService.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return &Claims{}, err
	}

	if !token.Valid {
		return &Claims{}, errors.New("invalid token")
	}

	return claims, nil
}

func (jwtService *JwtService) SignJWT(userID int) (string, error) {
	secret := []byte(jwtService.Secret)

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
