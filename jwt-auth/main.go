package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key
var jwtSecret = []byte("my-super-secret-key")

// Login request structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWT Claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

// LoginHandler handles POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Only POST request allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Receive JSON body
	var login LoginRequest

	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Normally these will come from database
	if login.Username != "jami" || login.Password != "123456" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create JWT claims
	claims := Claims{
		UserID:   1,
		Username: login.Username,
		Role:     "user",

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "my-auth-server",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create JWT token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign token using secret key
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   signedToken,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func main() {

	http.HandleFunc("/login", LoginHandler)

	fmt.Println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}