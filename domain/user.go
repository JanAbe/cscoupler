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
}

// UserRepository interface
type UserRepository interface {
	NextID() string
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
}

// NewUser creates a new user or returns an error when the hashing of the password fails
func NewUser(email, password string) (*User, error) {
	id := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, errors.New("Error hashing password")
	}

	return &User{id, email, string(hash)}, nil
}

// CheckEmail checks if the email of the user
// is valid or not.
func (u *User) CheckEmail() bool {
	return strings.Contains(u.Email, "@")
}
