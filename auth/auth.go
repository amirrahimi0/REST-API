package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("monkey")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// LoginUser handles user login
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.Credentials true "User credentials"
// @Success 200 {string} string "Login successful"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedCreds Credentials
	err = db.QueryRow("SELECT username, password FROM users WHERE username = ?", creds.Username).Scan(&storedCreds.Username, &storedCreds.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// LoginBookkeeper handles bookkeeper login
// @Summary Bookkeeper login
// @Description Authenticate a bookkeeper and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.Credentials true "Bookkeeper credentials"
// @Success 200 {string} string "Login successful"
// @Failure 401 {string} string "Unauthorized"
// @Router /login/bookkeepers [post]
func LoginBookkeeper(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "/Users/amir/Documents/newtestgo/test.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var storedCreds Credentials
	err = db.QueryRow("SELECT email, password FROM users WHERE email = ? AND role = 'admin'", creds.Username).Scan(&storedCreds.Username, &storedCreds.Password)
	if err != nil {
		fmt.Fprintf(w, "email: %s ", creds.Username)
		http.Error(w, "Bookkeeper not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// AuthMiddleware is a middleware for authenticating users
// @Summary User authentication middleware
// @Description Middleware to authenticate users using JWT token
// @Tags auth
// @Produce json
// @Success 200 {string} string "Authenticated"
// @Failure 401 {string} string "Unauthorized"
// @Router /user [get]
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// BookkeeperMiddleware is a middleware for authenticating bookkeepers
// @Summary Bookkeeper authentication middleware
// @Description Middleware to authenticate bookkeepers using JWT token
// @Tags auth
// @Produce json
// @Success 200 {string} string "Authenticated"
// @Failure 401 {string} string "Unauthorized"
// @Router /admin [get]
func BookkeeperMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminHandler handles requests to the admin page
// @Summary Admin page
// @Description Admin page accessible only to authenticated bookkeepers
// @Tags auth
// @Produce plain
// @Success 200 {string} string "Welcome to the admin page!"
// @Router /admin [get]
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the admin page!"))
}

// UserHandler handles requests to the user page
// @Summary User page
// @Description User page accessible only to authenticated users
// @Tags auth
// @Produce plain
// @Success 200 {string} string "Welcome to the user page!"
// @Router /user [get]
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the user page!"))
}
