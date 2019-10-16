package handlers

// This file contains all handlers regarding
// authorization and authentication,
// such as ValidateRequests and Signin.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/janabe/cscoupler/services"

	"github.com/dgrijalva/jwt-go"
)

// AuthHandler struct containing all authorization
// related handler/middleware funcs
type AuthHandler struct {
	JWTKey      []byte
	UserService services.UserService
}

// UserData is a struct that corresponds to incoming user data
type UserData struct {
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Claims is a struct to convey the second part of the JWT (sometimes called payload)
type Claims struct {
	ID     string
	Email  string
	UserID string
	jwt.StandardClaims
}

// Signin returns a handler for signin requests, creating
// a JWT for the user if all credentials are correct
// and storing this token in a cookie
func (a AuthHandler) Signin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}

		var data UserData

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if account with email exists
		user, err := a.UserService.FindByEmail(strings.ToLower(data.Email))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isValid := a.UserService.ValidatePassword(user.HashedPassword, data.Password)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		roleID, err := a.UserService.FindRoleID(user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Build the claims part of the JWT and
		// set the expiration time of the JWT (todo: find out what a good time is)
		expirationTime := time.Now().Add(6 * time.Hour)
		claims := &Claims{
			ID:     roleID,
			Email:  user.Email,
			UserID: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Create new token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(a.JWTKey)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			// todo: fix -> doesn't work with httponly and secure flags enabled
		})
	})
}

// Validate returns a handler used to secure endpoints.
// It validates incoming requests by checking if the user has a valid
// token and the correct role, and is thus allowed to call this endpoint or not.
// If the token is valid, h.serveHTTP() gets called which means the page is shown.
// If the role param is left empty (""), all roles are allowed to call this endpoint.
func (a AuthHandler) Validate(role string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := a.GetToken(cookie)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userEmail, ok := claims["Email"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := a.UserService.FindByEmail(userEmail)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if role != "" {
			if user.Role != role {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// If jwt is valid, serve the webpage of h. Aka run handler h
		h.ServeHTTP(w, r)
	})
}

// Register registers all authentication related handlers
func (a AuthHandler) Register() {
	http.Handle("/signin", LoggingHandler(os.Stdout, a.Signin()))
}

// GetToken gets the token from the cookie
func (a AuthHandler) GetToken(cookie *http.Cookie) (*jwt.Token, error) {
	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Incorrect Signing method used: %v", token.Header["alg"])
		}

		return a.JWTKey, nil
	})

	return token, err
}
