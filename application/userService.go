package application

import (
	"github.com/janabe/cscoupler/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserService struct, containing all features
// the app supports regaring users
type UserService struct {
	UserRepo domain.UserRepository
}

// Create creates a new Student
func (u UserService) Create(user *domain.User) error {
	err := u.UserRepo.Create(user)
	return err
}

// FindByID finds a user based on id
func (u UserService) FindByID(id string) (*domain.User, error) {
	user, err := u.UserRepo.FindByID(id)
	return user, err
}

// FindByEmail finds a user based on email
func (u UserService) FindByEmail(email string) (*domain.User, error) {
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
