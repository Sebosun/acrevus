package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"sebosun/acrevus-go/internal/database"

	"github.com/gin-gonic/gin"
)

func userTestRouter(config *ApiConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config.registerUserRoutes(router.Group("/api/v1"))
	return router
}

func userRequest(t *testing.T, router http.Handler, method, path, body string, status int) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "/api/v1"+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != status {
		t.Fatalf("%s %s: got %d, want %d: %s", method, path, rec.Code, status, rec.Body.String())
	}
	return rec
}

func responseUser(t *testing.T, rec *httptest.ResponseRecorder) UserResponse {
	t.Helper()
	if strings.Contains(strings.ToLower(rec.Body.String()), "password") {
		t.Fatal("response includes a password field")
	}
	var body struct {
		User UserResponse `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.User
}

func TestUserValidation(t *testing.T) {
	// A nil DB ensures invalid requests are rejected before database access.
	router := userTestRouter(&ApiConfig{})
	for _, tc := range []struct {
		name, method, path, body string
	}{
		{"invalid ID", "GET", "/users/nope", ""},
		{"zero ID", "DELETE", "/users/0", ""},
		{"negative ID", "PATCH", "/users/-1", `{"name":"Updated"}`},
		{"overflow ID", "GET", "/users/9223372036854775808", ""},
		{"missing payload", "POST", "/users", ""},
		{"malformed JSON", "POST", "/users", "{"},
		{"missing password", "POST", "/users", `{"name":"Alice","email":"alice@example.com"}`},
		{"invalid email", "POST", "/users", `{"name":"Alice","email":"invalid","password":"secret"}`},
		{"long password", "POST", "/users", `{"name":"Alice","email":"alice@example.com","password":"` + strings.Repeat("a", maxPasswordLength+1) + `"}`},
		{"empty patch", "PATCH", "/users/1", `{}`},
		{"null patch", "PATCH", "/users/1", `null`},
		{"null fields", "PATCH", "/users/1", `{"name":null,"email":null,"password":null}`},
		{"unknown patch field", "PATCH", "/users/1", `{"unknown":"value"}`},
		{"empty name", "PATCH", "/users/1", `{"name":""}`},
		{"empty email", "PATCH", "/users/1", `{"email":""}`},
		{"invalid patch email", "PATCH", "/users/1", `{"email":"invalid"}`},
		{"empty password", "PATCH", "/users/1", `{"password":""}`},
		{"long patch password", "PATCH", "/users/1", `{"password":"` + strings.Repeat("a", maxPasswordLength+1) + `"}`},
		{"zero limit", "GET", "/users?limit=0", ""},
		{"large limit", "GET", "/users?limit=101", ""},
		{"invalid limit", "GET", "/users?limit=abc", ""},
		{"negative offset", "GET", "/users?offset=-1", ""},
		{"overflow offset", "GET", "/users?offset=2147483648", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := userRequest(t, router, tc.method, tc.path, tc.body, http.StatusBadRequest)
			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Fatalf("missing JSON error: %s", rec.Body.String())
			}
		})
	}
}

func TestUserLifecycle(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to a migrated Postgres database for the user lifecycle test")
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	// Keep the temporary table on the same session for every handler/query.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, "CREATE TEMP TABLE users (LIKE public.users INCLUDING ALL)"); err != nil {
		t.Fatal(err)
	}
	queries := database.New(db)
	router := userTestRouter(&ApiConfig{DB: queries})
	call := func(method, path, body string, status int) *httptest.ResponseRecorder {
		return userRequest(t, router, method, path, body, status)
	}
	list := func(path string) []UserResponse {
		rec := call("GET", path, "", http.StatusOK)
		if strings.Contains(strings.ToLower(rec.Body.String()), "password") {
			t.Fatal("list response includes a password field")
		}
		var body struct {
			Users []UserResponse `json:"users"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Users == nil {
			t.Fatal("expected an array, not null")
		}
		return body.Users
	}
	if len(list("/users")) != 0 {
		t.Fatal("expected empty users list")
	}
	payload := `{"name":"Alice","email":"alice@example.com","password":"ExamplePass123!"}`
	created := responseUser(t, call("POST", "/users", payload, http.StatusOK))
	if created.ID <= 0 || created.Name != "Alice" || created.Email != "alice@example.com" ||
		created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() || created.DeletedAt != nil {
		t.Fatalf("unexpected created user: %+v", created)
	}
	stored, err := queries.GetUser(ctx, created.ID)
	if err != nil || !CheckPassword("ExamplePass123!", stored.Password) {
		t.Fatalf("password was not correctly hashed: %v", err)
	}
	path := fmt.Sprintf("/users/%d", created.ID)
	fetched := responseUser(t, call("GET", path, "", http.StatusOK))
	if fetched.ID != created.ID || fetched.Email != created.Email {
		t.Fatalf("unexpected fetched user: %+v", fetched)
	}
	call("POST", "/users", payload, http.StatusConflict)
	second := responseUser(t, call("POST", "/users",
		`{"name":"Bob","email":"bob@example.com","password":"AnotherPass123!"}`, http.StatusOK))
	page := list("/users?limit=1&offset=1")
	if len(page) != 1 || page[0].ID != second.ID {
		t.Fatalf("unexpected paginated users: %+v", page)
	}
	updated := responseUser(t, call("PATCH", path, `{"name":"Updated Alice"}`, http.StatusOK))
	afterUpdate, err := queries.GetUser(ctx, created.ID)
	if err != nil || updated.Name != "Updated Alice" || updated.Email != created.Email ||
		afterUpdate.Password != stored.Password || !updated.CreatedAt.Equal(created.CreatedAt) ||
		!updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatalf("partial update did not preserve omitted fields/timestamps: %+v, %v", updated, err)
	}
	call("PATCH", path, `{"email":"bob@example.com"}`, http.StatusConflict)
	responseUser(t, call("PATCH", path, `{"email":"alice.updated@example.com","password":"NewPass123!"}`, http.StatusOK))
	afterPassword, err := queries.GetUser(ctx, created.ID)
	if err != nil || afterPassword.Email != "alice.updated@example.com" ||
		!CheckPassword("NewPass123!", afterPassword.Password) || CheckPassword("ExamplePass123!", afterPassword.Password) {
		t.Fatalf("email/password update failed: %v", err)
	}
	if rec := call("DELETE", path, "", http.StatusNoContent); rec.Body.Len() != 0 {
		t.Fatal("204 response must have no body")
	}
	var deletedAt, updatedAt time.Time
	if err := db.QueryRowContext(ctx, "SELECT deleted_at, updated_at FROM users WHERE id = $1", created.ID).Scan(&deletedAt, &updatedAt); err != nil {
		t.Fatal(err)
	}
	if deletedAt.IsZero() || !deletedAt.Equal(updatedAt) {
		t.Fatal("soft deletion did not maintain timestamps")
	}
	call("GET", path, "", http.StatusNotFound)
	call("PATCH", path, `{"name":"Cannot restore"}`, http.StatusNotFound)
	call("DELETE", path, "", http.StatusNotFound)
	call("GET", "/users/9223372036854775807", "", http.StatusNotFound)
	remaining := list("/users")
	if len(remaining) != 1 || remaining[0].ID != second.ID {
		t.Fatalf("soft-deleted user still listed: %+v", remaining)
	}
	if _, err := queries.GetUserByEmail(ctx, afterPassword.Email); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("soft-deleted user still accessible by email: %v", err)
	}
	call("POST", "/users", `{"name":"Reuse","email":"alice.updated@example.com","password":"ExamplePass123!"}`, http.StatusConflict)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	rec := call("GET", "/users", "", http.StatusInternalServerError)
	if strings.Contains(rec.Body.String(), "database is closed") {
		t.Fatal("internal database error leaked to client")
	}
}
