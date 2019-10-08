package handlers

// HandlerFunc creates a handler from a normal func

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/janabe/cscoupler/application"
	"github.com/janabe/cscoupler/database/memory"
	"github.com/janabe/cscoupler/domain"
	"github.com/janabe/cscoupler/util"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = util.GetJWTSecret("./.secret.json")
var userRepo = memory.UserRepo{DB: make(map[string]domain.User)}
var userService = application.UserService{UserRepo: userRepo}

// Useronly ...
var Useronly = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello user"))
})

// SignupHandler is a handler for signup requests, creating a new
// user with the provided credentials
var SignupHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var creds Credentials

	// check if json is invalid
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userService.EmailAlreadyUsed(creds.Email) {
		fmt.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	user, err := domain.NewUser(creds.Email, creds.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = userService.Create(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
})

// SigninHandler is a handler for signin requests, creating
// a JWT for the user if all credentials are correct
// and storing this token in a cookie
// -----------------
// todo: What i don't understand is how are the responseWriter and Request
// passed to the signinHandler?
var SigninHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if account with email exists
	user, err := userService.FindByEmail(creds.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isValid := userService.ValidatePassword(user.HashedPassword, creds.Password)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Build the claims part of the JWT and
	// set the expiration time of the JWT (todo: find out what a good time is)
	expirationTime := time.Now().Add(6 * time.Hour)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		// doesn't work with httponly and secure flags enabled
	})
})

// ValidateHandler is a handler/middleware used to secure endpoints.
// It validates incoming requests by checking if the user has a valid
// token and is thus allowed to call this endpoint or not.
// If the token is valid, h.serveHTTP() gets called which means the page is thus shown.
func ValidateHandler(h http.Handler) http.Handler {
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

		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Incorrect Signing method used: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})

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

		// If jwt is valid, serve the webpage of h
		h.ServeHTTP(w, r)
	})
}

// ------------- HELP Structs -------------

// Credentials is a struct that maps to the credentials part of the request body
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims is a struct to convey the second part of the JWT (sometimes called payload)
type Claims struct {
	Email string
	jwt.StandardClaims
}
