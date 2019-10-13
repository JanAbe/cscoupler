package memory

import (
	"errors"

	"github.com/janabe/cscoupler/domain"
)

// UserRepo ...
type UserRepo struct {
	DB map[string]domain.User
}

// Create ...
func (u UserRepo) Create(user domain.User) error {
	u.DB[user.ID] = user
	return nil
}

// FindByID ...
func (u UserRepo) FindByID(id string) (domain.User, error) {
	if user, ok := u.DB[id]; ok {
		return user, nil
	}

	return domain.User{}, errors.New("No user with id " + id)
}

// FindByEmail ...
func (u UserRepo) FindByEmail(email string) (domain.User, error) {
	for _, user := range u.DB {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, errors.New("No user with email " + email)
}
