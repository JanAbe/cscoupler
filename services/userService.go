package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/janabe/cscoupler/domain"
	e "github.com/janabe/cscoupler/errors"
)

// UserService struct, containing all features
// the app supports regaring users
type UserService struct {
	UserRepo domain.UserRepository
}

// Register registers a user
func (u UserService) Register(user domain.User) error {
	if u.EmailAlreadyUsed(user.Email) {
		return e.ErrorEmailAlreadyUsed
	}

	err := u.UserRepo.Create(user)
	return err
}

// FindByEmail finds a user based on email
func (u UserService) FindByEmail(email string) (domain.User, error) {
	user, err := u.UserRepo.FindByEmail(email)
	return user, err
}

// EmailAlreadyUsed checks if the email is already used
// for an account in the system
func (u UserService) EmailAlreadyUsed(email string) bool {
	_, err := u.FindByEmail(email)
	if err != nil {
		return false
	}

	return true
}

// ValidatePassword validates a hashed password with an unhashed password
// returning true if the hashed password is the hash of the provided password
// returning false otherwise
func (u UserService) ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
