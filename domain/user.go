package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User struct to convey user
type User struct {
	ID             string
	Email          string
	HashedPassword string
	Firstname      string
	Lastname       string
	Role           string
}

// UserRepository interface
type UserRepository interface {
	NextID() string
	Create(user User) error
	FindByID(id string) (User, error)
	FindByEmail(email string) (User, error)
}

// NewUser creates a new user or returns an error when the hashing of the password fails
func NewUser(email, password, fname, lname, role string) (User, error) {
	id := uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("Error hashing password")
	}

	if len(strings.TrimSpace(fname)) == 0 {
		return User{}, errors.New("provided firstname can't be empty")
	}

	if len(strings.TrimSpace(lname)) == 0 {
		return User{}, errors.New("provided lastname can't be empty")
	}

	return User{
		ID:             id,
		Email:          strings.ToLower(email),
		HashedPassword: string(hash),
		Firstname:      strings.ToLower(fname),
		Lastname:       strings.ToLower(lname),
		Role:           role,
	}, nil
}
