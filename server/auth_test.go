package server

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignAndVerifyJWT(t *testing.T) {
	config := APIConfig{
		ENVS: Envs{JwtSecret: "test-secret"},
	}

	signedToken, err := config.signJWT(42)
	if err != nil {
		t.Fatalf("signJWT returned an error: %v", err)
	}

	token, err := config.verifyJWT(signedToken)
	if err != nil {
		t.Fatalf("verifyJWT returned an error: %v", err)
	}

	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		t.Fatalf("unexpected signing method: %s", token.Method.Alg())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", token.Claims)
	}
	userID := claims["user_id"]
	if userID != float64(42) {
		t.Fatalf("unexpected user_id claim: %v", userID)
	}
}
