package server

import (
	"strings"
	"testing"
)

func TestArgon2PasswordHash(t *testing.T) {
	password := "ExamplePass123!"
	first, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	second, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("identical passwords must use unique salts")
	}
	if !strings.HasPrefix(first, "$argon2id$v=19$m=19456,t=2,p=1$") {
		t.Fatalf("unexpected Argon2id hash format: %q", first)
	}
	if !CheckPassword(password, first) {
		t.Fatal("correct password was rejected")
	}
	if CheckPassword("incorrect", first) {
		t.Fatal("incorrect password was accepted")
	}
}

func TestArgon2PasswordHashRejectsInvalidInput(t *testing.T) {
	if _, err := HashPassword(strings.Repeat("a", maxPasswordLength+1)); err != ErrPasswordTooLong {
		t.Fatalf("HashPassword error = %v, want ErrPasswordTooLong", err)
	}
	for _, hash := range []string{
		"",
		"$bcrypt$...",
		"$argon2id$v=18$m=19456,t=2,p=1$c2FsdHNhbHQ$MTIzNDU2Nzg5MDEyMzQ1",
		"$argon2id$v=19$m=999999,t=2,p=1$c2FsdHNhbHQ$MTIzNDU2Nzg5MDEyMzQ1",
	} {
		if CheckPassword("password", hash) {
			t.Fatalf("invalid hash accepted: %q", hash)
		}
	}
}
